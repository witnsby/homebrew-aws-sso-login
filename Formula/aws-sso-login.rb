class AwsSsoLogin < Formula
  desc "CLI that streamlines AWS SSO authentication and credentials management"
  homepage "https://github.com/witnsby/aws-sso-login"
  version "0.0.7"
  license "Apache-2.0"

  on_macos do
    on_intel do
      url "https://github.com/witnsby/aws-sso-login/releases/downloadsso-login_darwin_amd64"
      sha256 "f67fd7f9ac6d189d4"
    end
    on_arm do
      url "https://github.com/witnsby/aws-sso-log7/aws-sso-login_darwin_arm64"
      sha256 "1a03ca7a05c92da418b6f9358d56ac213db44eb7a5ad7a2cc35ddac1a7ceb7c0"
    end
  end

  on_linux do
    on_intel do
      url "https://github.com/witnsby/aws-sso-login/releases/download/v0.0.7/aws-sso-login_linux_amd64"
      sha256 "bb71e8171b5fd3db84939f90a023827c44b6a0a4ca12edad09291363147ab84f"
    end
    on_arm do
      url "https://github.com/witnsby/aws-sso-login/releases/download/v0.0.7/aws-sso-login_linux_arm64"
      sha256 "0aea07e02203919828bbbaae17198b9439a3bf2dec5c7175a553c83150fa8230"
    end
  end

  def install
    bin.install File.basename(stable.url) => "aws-sso-login"
  end

  test do
    assert_match version.to_s, shell_output("#{bin}/aws-sso-login version")
  end
end
