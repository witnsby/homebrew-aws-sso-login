class AwsSsoLogin < Formula
  desc "CLI that streamlines AWS SSO authentication and credentials management"
  homepage "https://github.com/witnsby/aws-sso-login"
  version "0.2.1"
  license "Apache-2.0"

  on_macos do
    on_intel do
      url "https://github.com/witnsby/aws-sso-login/releases/download/v0.2.1/aws-sso-login_darwin_amd64"
      sha256 "c49b80922fbd5c06adbfc2b4b9b80a6280729ce061a9c5bdec9d44f0e363897e"
    end
    on_arm do
      url "https://github.com/witnsby/aws-sso-login/releases/download/v0.2.1/aws-sso-login_darwin_arm64"
      sha256 "0bd0da27ed8f9584f7df9dbf3a5bba85df4607be549df9c5ba67f5ab469d4ed8"
    end
  end

  on_linux do
    on_intel do
      url "https://github.com/witnsby/aws-sso-login/releases/download/v0.2.1/aws-sso-login_linux_amd64"
      sha256 "50dc805173c97fa1de82230c63a6e912f99c26b0bf5e8bb3a5b0fccc963b3b2b"
    end
    on_arm do
      url "https://github.com/witnsby/aws-sso-login/releases/download/v0.2.1/aws-sso-login_linux_arm64"
      sha256 "856202e005fc5cf62c47f125eb3b1cec043cb09fecdfe48f3b437900f0e9f5f1"
    end
  end

  def install
    bin.install File.basename(stable.url) => "aws-sso-login"
  end

  test do
    assert_match version.to_s, shell_output("#{bin}/aws-sso-login version")
  end
end
