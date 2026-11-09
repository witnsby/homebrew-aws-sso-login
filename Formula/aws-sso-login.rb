class AwsSsoLogin < Formula
  desc "CLI that streamlines AWS SSO authentication and credentials management"
  homepage "https://github.com/witnsby/aws-sso-login"
  version "0.3.0"
  license "Apache-2.0"

  on_macos do
    on_intel do
      url "https://github.com/witnsby/aws-sso-login/releases/download/v0.3.0/aws-sso-login_darwin_amd64"
      sha256 "33ddd21955b6d65c6d8237012e5119615590979ce0b157159eb3a47000243c68"
    end
    on_arm do
      url "https://github.com/witnsby/aws-sso-login/releases/download/v0.3.0/aws-sso-login_darwin_arm64"
      sha256 "4ddf7e4e6f4f9bae4db2cc9621af9f89edfd7fd3edb3bcc79e1a96d44537cb26"
    end
  end

  on_linux do
    on_intel do
      url "https://github.com/witnsby/aws-sso-login/releases/download/v0.3.0/aws-sso-login_linux_amd64"
      sha256 "46053fdb4e2c6fcf0517a232c85df950691bf4342aa46bc7716b1c1cf0dfd3ec"
    end
    on_arm do
      url "https://github.com/witnsby/aws-sso-login/releases/download/v0.3.0/aws-sso-login_linux_arm64"
      sha256 "e80b4e54c0169cefd80c5360955cf89026f1e7c85d675eed8865ec58af631baa"
    end
  end

  def install
    bin.install File.basename(stable.url) => "aws-sso-login"
  end

  test do
    assert_match version.to_s, shell_output("#{bin}/aws-sso-login version")
  end
end
