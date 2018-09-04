#debuginfo not supported with Go
%global debug_package %{nil}

Name: argo-api-authn
Summary: ARGO Authentication API. Map X509, OICD to token.
Version: 0.1.0
Release: 1%{?dist}
License: ASL 2.0
Buildroot: %{_tmppath}/%{name}-buildroot
Group: Unspecified
Source0: %{name}-%{version}.tar.gz
BuildRequires: golang
BuildRequires: git
Requires(pre): /usr/sbin/useradd, /usr/bin/getent
ExcludeArch: i386

%description
Installs the ARGO Authentication API

%pre
/usr/bin/getent group argo-api-authn || /usr/sbin/groupadd -r argo-api-authn
/usr/bin/getent passwd argo-api-authn || /usr/sbin/useradd -r -s /sbin/nologin -d /var/www/argo-api-authn -g argo-api-authn argo-api-authn

%prep
%setup

%build
export GOPATH=$PWD
export PATH=$PATH:$GOPATH/bin

cd src/github.com/ARGOeu/argo-api-authn/
go install

%install
%{__rm} -rf %{buildroot}
install --directory %{buildroot}/var/www/argo-api-authn
install --mode 755 bin/argo-api-authn %{buildroot}/var/www/argo-api-authn/argo-api-authn

install --directory %{buildroot}/etc/argo-api-authn
install --directory %{buildroot}/etc/argo-api-authn/conf.d
install --mode 644 src/github.com/ARGOeu/argo-api-authn/conf/argo-api-authn-config.template %{buildroot}/etc/argo-api-authn/conf.d/argo-api-authn-config.json

install --directory %{buildroot}/usr/lib/systemd/system
install --mode 644 src/github.com/ARGOeu/argo-api-authn/argo-api-authn.service %{buildroot}/usr/lib/systemd/system/

%clean
%{__rm} -rf %{buildroot}
export GOPATH=$PWD
cd src/github.com/ARGOeu/argo-api-authn/
go clean

%files
%defattr(0644,argo-api-authn,argo-api-authn)
%attr(0750,argo-api-authn,argo-api-authn) /var/www/argo-api-authn
%attr(0755,argo-api-authn,argo-api-authn) /var/www/argo-api-authn/argo-api-authn
%config(noreplace) %attr(0644,argo-api-authn,argo-api-authn) /etc/argo-api-authn/conf.d/argo-api-authn-config.json
%attr(0644,root,root) /usr/lib/systemd/system/argo-api-authn.service

%changelog
* Thu Jun 14 2018 Themis Zamani  <themiszamani@gmail.com> - 0.1.0-1%{?dist}
- Initial release
