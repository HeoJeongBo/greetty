# greetty zsh hook — sourced once per interactive session.
# Shows the greeting only on the first prompt so it doesn't repeat.
if [[ -z "$GREETTY_SHOWN" ]]; then
  export GREETTY_SHOWN=1
  command greetty greet 2>/dev/null
fi
