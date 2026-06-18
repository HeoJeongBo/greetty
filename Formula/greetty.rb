class Greetty < Formula
  desc "Pretty, developer-flavored greeting banner for your terminal"
  homepage "https://github.com/HeoJeongBo/greetty"
  url "https://github.com/HeoJeongBo/greetty/archive/refs/tags/v0.1.0.tar.gz"
  sha256 "fffa2d37513f39c36d061e754b59ef068472857e3e204f629c9edb8cb806b69d"
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
