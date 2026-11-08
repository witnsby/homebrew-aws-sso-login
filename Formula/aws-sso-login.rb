class AwsSsoLogin < Formula
  desc "CLI that streamlines AWS SSO authentication and credentials management"
  homepage "https://github.com/witnsby/aws-sso-login"
  version "0.2.0"
  license "Apache-2.0"

  on_macos do
    on_intel do
      url "https://github.com/witnsby/aws-sso-login/releases/download/v0.2.0/aws-sso-login_darwin_amd64"
      sha256 "2ca8cb6b02544d19e62c9f6fc2d98dde2b456941d1499120c19218cedc4adeaf"
    end
    on_arm do
      url "https://github.com/witnsby/aws-sso-login/releases/download/v0.2.0/aws-sso-login_darwin_arm64"
      sha256 "0f71ada23de3edc1e5d080e0ecc5b338e034954dd4abebf60473a827c28445f9"
    end
  end

  on_linux do
    on_intel do
      url "https://github.com/witnsby/aws-sso-login/releases/download/v0.2.0/aws-sso-login_linux_amd64"
      sha256 "12ebdd18c45ceb34ba95d8843cdded0001cc6b7a7430ec4da6393f8eb9dda6cd"
    end
    on_arm do
      url "https://github.com/witnsby/aws-sso-login/releases/download/v0.2.0/aws-sso-login_linux_arm64"
      sha256 "e59473cdf72ab4bdc498072694c80a74e8db797fe4223639798d423078e6f4e3"
    end
  end

  def install
    bin.install File.basename(stable.url) => "aws-sso-login"
  end

  test do
    assert_match version.to_s, shell_output("#{bin}/aws-sso-login version")
  end
end
