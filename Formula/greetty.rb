class Greetty < Formula
  desc "Pretty, developer-flavored greeting banner for your terminal"
  homepage "https://github.com/HeoJeongBo/greetty"
  url "https://github.com/HeoJeongBo/greetty/archive/refs/tags/v0.3.0.tar.gz"
  sha256 "ae7fa04330599c80625eb746d75e8bbc026f3f20b2b3954c6b348261122bf648"
  license "MIT"
  head "https://github.com/HeoJeongBo/greetty.git", branch: "main"

  depends_on "go" => :build

  def install
    ldflags = "-s -w -X github.com/HeoJeongBo/greetty/internal/cli.version=#{version}"
    system "go", "build", *std_go_args(ldflags: ldflags), "./cmd/greetty"
  end

  test do
    # greet must render the banner without touching the shell config
    assert_match "·", shell_output("#{bin}/greetty greet")
    assert_match version.to_s, shell_output("#{bin}/greetty --version")
  end
end
