# cfmagic

Clang is cool and Clang-Format is really cool, so adding genetic algorithms makes it super cool __and__ hype.

[![asciicast](https://asciinema.org/a/232250.svg)](https://asciinema.org/a/232250)

## Why ?

Clang-Format is a great tool when you want to make sure your team uses the proper indent conventions so that your code base looks homogenous.
Unfortunately, you will soon realize that there are plenty of clang-format configuration parameters, which will frequently interact with each other.

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

  - `<clang format path>`: (ex.: /usr/bin/clang-format-6.0) the path of your locally installed clang-format binary
  - `<population size>`: (ex.: 20) the size of the population (ie.: the number of randomly generated config files). A number between 10 and 100 is nice.
  - `<mutation rate>`: (ex: 14) the mutation rate. A number between 1 and 20 is nice

## Credits

The __configs.json__ file has been taken from [https://github.com/zed0/clang-format-configurator](https://github.com/zed0/clang-format-configurator), all credits to him.

## FAQ

__How long does it take to find a nice .clang-format configuration file matching my crazy indentation style?__

A few minutes is usually enough.

__It's been running for three weeks, how do I get the results now?__

Just hit ctrl+C, the current best solution will be written in a .clang-format file.

__Will cfmagic always find the ultimate .clang-format configuration file?__

No, simply because all possibilities are not managed by clang-format. cfmagic will find a "close" solution.

In addition to this, genetic algorithm do not necessarily find the best possible solution.

__What about tuning population size and mutation rate parameters?__

I find out that using a large population (ex. 100) seems to be counter productive.
Mutation rate can be tuned, but adding too much randomness (> 20) doesn't make any sense.

__I found your mutation boost secret feature!__

After a certain number of generations (ie. iterations), you may end up with a set of individuals (ie. clang-format configs) very close to each other.
At this point, you will want to renew things by adding more randomness.

In order to do so, we compute the standard deviation of the scores for every generation. If too low, mutationRate is doubled for the next generation.

__I see "signal: floating point exception" errors__

Funnily enough, cfmagic also acts as a clang-format fuzzer!
So "yes", some versions of clang-format may crash: this is handled by cfmagic and has no impact on the results (you may safely ignore these errors).

