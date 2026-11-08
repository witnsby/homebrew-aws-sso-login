class AwsSsoLogin < Formula
  desc "CLI that streamlines AWS SSO authentication and credentials management"
  homepage "https://github.com/witnsby/aws-sso-login"
  version "0.1.1"
  license "Apache-2.0"

  on_macos do
    on_intel do
      url "https://github.com/witnsby/aws-sso-login/releases/download/v0.1.1/aws-sso-login_darwin_amd64"
      sha256 "9f936ffe04372aabd63884045cf33ca90d741c189b8bc01259bb67f421fa20fd"
    end
    on_arm do
      url "https://github.com/witnsby/aws-sso-login/releases/download/v0.1.1/aws-sso-login_darwin_arm64"
      sha256 "20ce4ddc0a4978d9964261789e2da63a7c2533737687bb28b6b19aaa667f6000"
    end
  end

  on_linux do
    on_intel do
      url "https://github.com/witnsby/aws-sso-login/releases/download/v0.1.1/aws-sso-login_linux_amd64"
      sha256 "087e2554920eab5a5b576b9da4d37c18a21140be12b4d8ffdfb2e3fa7a2ea0db"
    end
    on_arm do
      url "https://github.com/witnsby/aws-sso-login/releases/download/v0.1.1/aws-sso-login_linux_arm64"
      sha256 "aa357deb10a6e9c47052e91683cd601a26d222a1034978adf1b0a85ee51d300d"
    end
  end

  def install
    bin.install File.basename(stable.url) => "aws-sso-login"
  end

  test do
    assert_match version.to_s, shell_output("#{bin}/aws-sso-login version")
  end
end
