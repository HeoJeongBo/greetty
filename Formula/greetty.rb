class Greetty < Formula
  desc "Pretty, developer-flavored greeting banner for your terminal"
  homepage "https://github.com/HeoJeongBo/greetty"
  url "https://github.com/HeoJeongBo/greetty/archive/refs/tags/v0.2.0.tar.gz"
  sha256 "1639a4abde490961c6d0b373502380531c64076e3cf08bb4cca40d32cee5c318"
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
