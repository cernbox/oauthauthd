# 
# oauthauthd spec file
#

Name: oauthauthd
Summary: Authentication daemon for CERNBox OCM implementation.
Version: 0.0.1
Release: 1%{?dist}
License: AGPLv3
BuildRoot: %{_tmppath}/%{name}-buildroot
Group: CERN-IT/ST
BuildArch: x86_64
Source: %{name}-%{version}.tar.gz

%description
This RPM provides a golang webserver that provides an authentication service for web clients.

# Don't do any post-install weirdness, especially compiling .py files
%define __os_install_post %{nil}

%prep
%setup -n %{name}-%{version}

%install
# server versioning

# installation
rm -rf %buildroot/
mkdir -p %buildroot/usr/local/bin
mkdir -p %buildroot/etc/oauthauthd
mkdir -p %buildroot/etc/logrotate.d
mkdir -p %buildroot/usr/lib/systemd/system
mkdir -p %buildroot/var/log/oauthauthd
install -m 755 oauthauthd	     %buildroot/usr/local/bin/oauthauthd
install -m 644 oauthauthd.service    %buildroot/usr/lib/systemd/system/oauthauthd.service
install -m 644 oauthauthd.yaml       %buildroot/etc/oauthauthd/oauthauthd.yaml
install -m 644 oauthauthd.logrotate  %buildroot/etc/logrotate.d/oauthauthd

%clean
rm -rf %buildroot/

%preun

%post

%files
%defattr(-,root,root,-)
/etc/oauthauthd
/etc/logrotate.d/oauthauthd
/var/log/oauthauthd
/usr/lib/systemd/system/oauthauthd.service
/usr/local/bin/*
%config(noreplace) /etc/oauthauthd/oauthauthd.yaml


%changelog
* Wed Oct 10 2018 Diogo Castro <diogo.castro@cern.ch> 0.0.1
- v0.0.1

