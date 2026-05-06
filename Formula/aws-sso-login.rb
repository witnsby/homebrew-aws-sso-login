class AwsSsoLogin < Formula
  desc "CLI that streamlines AWS SSO authentication and credentials management"
  homepage "https://github.com/witnsby/aws-sso-login"
  version "0.1.0"
  license "Apache-2.0"

  on_macos do
    on_intel do
      url "https://github.com/witnsby/aws-sso-login/releases/download/v0.1.0/aws-sso-login_darwin_amd64"
      sha256 "39d57b150db0452b3abe1d550effa976324d3ef208308e5713d991c2ebf67167"
    end
    on_arm do
      url "https://github.com/witnsby/aws-sso-login/releases/download/v0.1.0/aws-sso-login_darwin_arm64"
      sha256 "62a547f799e39e5bd3c68fc0b6cd9c1b1f0a04ec9af058d7376882c3b4a0b9d6"
    end
  end

  on_linux do
    on_intel do
      url "https://github.com/witnsby/aws-sso-login/releases/download/v0.1.0/aws-sso-login_linux_amd64"
      sha256 "2dccb92f20874931f9bb88951a0500cefc1b2b9e0dddd65d07721f67562280f4"
    end
    on_arm do
      url "https://github.com/witnsby/aws-sso-login/releases/download/v0.1.0/aws-sso-login_linux_arm64"
      sha256 "33753ebe8f90b619a7fbaf40712ff34238efb192ccf3e0e3eeb66e0effab6887"
    end
  end

  def install
    bin.install File.basename(stable.url) => "aws-sso-login"
  end

  test do
    assert_match version.to_s, shell_output("#{bin}/aws-sso-login version")
  end
end
