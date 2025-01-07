class AwsSsoLogin < Formula
  desc "CLI that streamlines AWS SSO authentication and credentials management"
  homepage "https://github.com/witnsby/aws-sso-login"
  url "https://github.com/witnsby/aws-sso-login/archive/refs/tags/v0.0.7.tar.gz"
  sha256 "ef6868d1285293fe06192c6125833f35595cd27e1a0bef9c5c5c74ff64ad675f"
  license "Apache-2.0"

  depends_on "go" => :build

  def install
    system "go", "build", *std_go_args(output: bin/"aws-sso-login"), "./src/cmd/bin/main.go"
  end

  test do
    assert_match version.to_s, shell_output("#{bin}/aws-sso-login version")
  end
end
