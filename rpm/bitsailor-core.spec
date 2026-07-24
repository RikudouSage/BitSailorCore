Name:       bitsailor-core
Version:    0.1.0
Release:    1
Summary:    BitSailor core library
License:    MIT

URL:        https://go.chrastecky.dev/bitsailor-core/bitwarden
Source0:    %{name}-%{version}.tar.gz

BuildRequires:  go
BuildRequires:  make
BuildRequires:  patchelf
BuildRequires:  gcc

%ifarch aarch64
%global make_target release-lib-arm64
%global out_dir arm64
%global cc_var CC_ARM64
%endif

%ifarch armv7hl
%global make_target release-lib-arm7
%global out_dir arm7
%global cc_var CC_ARMV7
%endif

%ifarch i486
%global make_target release-lib-386
%global out_dir 386
%global cc_var CC_386
%endif

%description
bitsailor-core provides the c-shared Bitwarden client library and C headers
used by Bitsailor

%prep
%setup -q

%build
%{!?make_target:%{error:Unsupported architecture %{_target_cpu}}}
export %{cc_var}="${CC:-gcc}"
make %{make_target}

%install
rm -rf %{buildroot}
install -d %{buildroot}%{_libdir}
install -d %{buildroot}%{_includedir}/bitsailor-core

install -m 0755 out/%{out_dir}/libbw.so %{buildroot}%{_libdir}/libbw.so
install -m 0644 out/%{out_dir}/*.h %{buildroot}%{_includedir}/bitsailor-core/

%files
%{_libdir}/libbw.so
%{_includedir}/bitsailor-core/*.h

