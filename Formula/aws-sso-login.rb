class AwsSsoLogin < Formula
  desc "CLI that streamlines AWS SSO authentication and credentials management"
  homepage "https://github.com/witnsby/aws-sso-login"
  url "https://github.com/witnsby/aws-sso-login/archive/refs/tags/v0.3.tar.gz"
  sha256 "b37b5262f4f680f"
  license "Apache-2.0"

  depends_on "go" => :build

  def install
    system "go", "build", *std_go_args(output: bin/"aws-sso-login"), "./src/cmd/bin/main.go"
  end

  test do
    assert_match version.to_s, shell_output("#{bin}/aws-sso-login version")
  end
end
