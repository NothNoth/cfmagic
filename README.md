# cfmagic

Clang is cool and Clang-Format is really cool, so adding genetic algorithms makes it super cool __and__ hype.

## Why ?

Clang-Format is a great tool when you want to make sure your team uses the proper indent conventions so that your code base looks homogenous.
Unfortunately, you will soon realize that there are plenty of clang-format configuration parameters, which usually will interact with each other.

You usually end up changing randomly some settings that might fix your problem for several hours.

When you have no real idea of how you should handle such a problem but the two following rules can be observed, genetic algorithms are a proper choice:

  - You can easily score a result
  - Mixing two good solutions may give you a better one


## How?

Basically cfmagic will apply the following algorithm:

  1. Generate a set of random .clang-format configuration files
  2. Apply clang-format with these config files on a piece of code properly indented according to your standards
  3. Compute the difference between the original piece of code and the re-formatted one and establish a 'score'
  4. Drop those with too many differences (high score)
  5. Mix the others to regenerate new ones (and add some randomness: mutations)
  6. GOTO 2


## Isn't it overkill?

Yes probably, but that's fun and actually is does work pretty well!

## How do I use this hype-overkill-wizardry?

__Build:__

```bash
  go build
```

Take a piece of code and format it manually according to your standards.
__Run:__

```bash
  ./cfmagic `<clang format path>` mypieceofcode.c `<population size>` `<mutation rate>`
```

You may want to adapt:

  - `<clang format path>`: (ex.: usr/bin/clang-format-6.0) the path of your locally installed clang-format binary
  - `<population size>`: (ex.: 20) the size of the population (ie.: the number of randomly generated config files). A number between 10 and 100 is nice.
  - `<mutation rate>`: (ex: 14) the mutation rate. A number between 1 and 20 is nice