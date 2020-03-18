package mysqluserbackend

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"github.com/bluele/gcache"
	"github.com/cernbox/oauthauthd/pkg"
	_ "github.com/go-sql-driver/mysql"
	"go.uber.org/zap"
	"time"
)

type Options struct {
	Hostname  string
	Port      int
	Username  string
	Password  string
	DB        string
	CacheSize int
	CacheTTL  int
	Logger    *zap.Logger
}

func New(opt *Options) pkg.UserBackend {
	setDefaults(opt)

	cache := gcache.New(opt.CacheSize).LFU().Build()

	db, err := sql.Open("mysql", fmt.Sprintf("%s:%s@tcp(%s:%d)/%s", opt.Username, opt.Password, opt.Hostname, opt.Port, opt.DB))
	if err != nil {
		panic(fmt.Sprintf("CANNOT CONNECT TO MYSQL SERVER. MAKE SURE THIS SERVICE RUNS AFTER DB IS UP AND RUNNING. HOSTNAME=%s PORT=%d DB=%s", opt.Hostname, opt.Port, opt.DB))
	}

	return &userBackend{
		db:       db,
		logger:   opt.Logger,
		cache:    cache,
		cacheTTL: time.Second * time.Duration(opt.CacheTTL),
	}
}

func setDefaults(opt *Options) {
	if opt.CacheSize == 0 {
		opt.CacheSize = 1000000 // 1 million
	}

	if opt.CacheTTL == 0 {
		opt.CacheTTL = 60 // seconds
	}
}

type userBackend struct {
	hostname string
	port     int
	username string
	password string
	db       *sql.DB
	table    string

	logger   *zap.Logger
	cache    gcache.Cache
	cacheTTL time.Duration
}

// TODO implement caching

// returns empty token if no entry is available.
func (ub *userBackend) getFromCache(ctx context.Context, key string) (ti *tokenInfo) {
	v, err := ub.cache.Get(key)
	if err == nil {
		if tokenInfo, ok := v.(*tokenInfo); ok {
			ti = tokenInfo
		}
	}
	return
}

// set key in cache with timet=-to-live in seconds
func (ub *userBackend) storeInCache(ctx context.Context, key string, ti *tokenInfo) {
	ub.cache.SetWithExpire(key, ti, ub.cacheTTL)
}

type tokenInfo struct {
	userId  string
	expires int64
}

func newTokenInfo(userId string, expires int64) *tokenInfo {
	return &tokenInfo{userId: userId, expires: expires}
}

func (ub *userBackend) Authenticate(ctx context.Context, token string) (string, error) {
	// check if token is still in the cache
	cachedTokenInfo := ub.getFromCache(ctx, token)
	if cachedTokenInfo != nil && cachedTokenInfo.userId != "" {
		ub.logger.Info("OAUTH AUTHENTICATED BY USING CACHE", zap.String("user", cachedTokenInfo.userId))
		return cachedTokenInfo.userId, nil
	}

	var user string
	var expires int64

	query := "SELECT user_id, expires FROM oc_oauth2_access_tokens WHERE token=?"
	err := ub.db.QueryRow(query, token).Scan(&user, &expires)
	if err != nil {
		if err == sql.ErrNoRows {
			return "", errors.New("Token not found")
		}
		ub.logger.Error("CANNOT QUERY STATEMENT")
		return "", err
	}
	now := time.Now().Unix()

	if expires < now {
		return "", errors.New("Expired")
	}

	ub.logger.Info("OAUTH AUTHENTICATED", zap.String("user", user))

	// store in the cache
	ti := newTokenInfo(user, expires)
	ub.storeInCache(ctx, token, ti)

	ub.logger.Info("OAUTH TOKEN STORED IN CACHE", zap.String("user", user), zap.Int64("expire", expires))

	return user, nil
}

func (ub *userBackend) SetExpiration(ctx context.Context, expiration int64) error {

	return nil
}

func (ub *userBackend) ClearCache(ctx context.Context) {

}
