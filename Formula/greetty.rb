class Greetty < Formula
  desc "Pretty, developer-flavored greeting banner for your terminal"
  homepage "https://github.com/HeoJeongBo/greetty"
  url "https://github.com/HeoJeongBo/greetty/archive/refs/tags/v0.1.0.tar.gz"
  # Replace with: shasum -a 256 v0.1.0.tar.gz   (see DEPLOY.md step 4)
  sha256 "REPLACE_WITH_TARBALL_SHA256"
  license "MIT"
  head "https://github.com/HeoJeongBo/greetty.git", branch: "main"

  depends_on "go" => :build

  def install
    ldflags = "-s -w -X github.com/HeoJeongBo/greetty/cmd.version=#{version}"
    system "go", "build", *std_go_args(ldflags: ldflags)
  end

  test do
    # greet must render the banner without touching the shell config
    assert_match "·", shell_output("#{bin}/greetty greet")
    assert_match version.to_s, shell_output("#{bin}/greetty --version")
  end
end
