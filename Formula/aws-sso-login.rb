class AwsSsoLogin < Formula
  desc "CLI that streamlines AWS SSO authentication and credentials management"
  homepage "https://github.com/witnsby/aws-sso-login"
  version "0.0.8"
  license "Apache-2.0"

  on_macos do
    on_intel do
      url "https://github.com/witnsby/aws-sso-login/releases/download/v0.0.8/aws-sso-login_darwin_amd64"
      sha256 "63e34dc09fbd244b798824a9921fb564ef6051f16cd041d2b58b654e1bad6ccf"
    end
    on_arm do
      url "https://github.com/witnsby/aws-sso-login/releases/download/v0.0.8/aws-sso-login_darwin_arm64"
      sha256 "f8fe0da63c90ef068f6e3964b426fc5681c596e1d99f8e6abbd429b97b21aa6c"
    end
  end

  on_linux do
    on_intel do
      url "https://github.com/witnsby/aws-sso-login/releases/download/v0.0.8/aws-sso-login_linux_amd64"
      sha256 "9663cce8b120cc38a097f223b27c273759c3a573a02bcf8036651417b1e45d14"
    end
    on_arm do
      url "https://github.com/witnsby/aws-sso-login/releases/download/v0.0.8/aws-sso-login_linux_arm64"
      sha256 "62741595d59f74c7f873abc534e32a32016520b8ab5604522de7025a01816677"
    end
  end

  def install
    bin.install File.basename(stable.url) => "aws-sso-login"
  end

  test do
    assert_match version.to_s, shell_output("#{bin}/aws-sso-login version")
  end
end
