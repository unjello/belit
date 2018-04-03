// andrzej lichnerowicz, unlicensed (~public domain)

int cause_a_warning() {}
int main() {
  return cause_a_warning();
}

/* belit: cxx=g++-7 cxxopts="-Wall -Werror" */