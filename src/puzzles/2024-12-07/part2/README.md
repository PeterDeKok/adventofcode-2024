# Advent of Code 2024 | Day 7 | Part 2

Eric Wastl deserves the credit for creating these challenges year after year. Therefore, the problem description, input and Solutions won't be available in this repo.
However, a link to the problem description will be provided and the solution can be generated locally.

| Year        | 2024                    |
|-------------|-------------------------|
| Date        | 2024-12-07              |
| Part        | 2                       |
| Finished on | 2024-12-07 11:40:27 CET |

## Problem

[adventofcode.com/2024/day/7](https://adventofcode.com/2024/day/7)

## Runtime

```
21.457967ms
```

## Notable events (post mortem)

This challenge had some interesting results. To me one of them still seems a bit surreal.

- Logging comes with a serious speed penalty
- Recursion is faster

This all started because a seemingly simple program was running for over 5 minutes...

### Logging comes with a serious speed penalty

Logging has had roughly a 2x speed penalty. To up to 50x in one case!!!

With `-tags=DEBUG`:
```
BenchmarkSolution_Sample-20                          	   33452	     35849 ns/op
BenchmarkSolution_Input-20                           	       2	 842071197 ns/op
BenchmarkParseLine-20                                	 2771608	       428.6 ns/op
BenchmarkHandleLine_bitwise-20                       	       1	17634369985 ns/op
BenchmarkHandleLine_bitwise_withoutFullLogging-20    	       3	 423928658 ns/op
BenchmarkHandleLine_recursion-20                     	      10	 104103530 ns/op
```

Without build tags:
```
BenchmarkSolution_Sample-20                          	   96205	     12779 ns/op
BenchmarkSolution_Input-20                           	     100	  13947979 ns/op
BenchmarkParseLine-20                                	 2795397	       431.7 ns/op
BenchmarkHandleLine_bitwise-20                       	       1	16770730541 ns/op
BenchmarkHandleLine_bitwise_withoutFullLogging-20    	       4	 289249480 ns/op
BenchmarkHandleLine_recursion-20                     	     811	   1367284 ns/op
```

While I never expected it to be free; nor cheap. I was pretty flaborgasted to see the 50x difference.

I have not looked too deep into it. I did add some magic to the logging so the TUI can show the logs if I ever finish it... There could easily be a huge blocking speed penalty to this code.
What I do know is that - while not the bulk of the time penalty - composing strings is a huge factor to this punishment.

But either way, I've now created a way to disable logging, which is also the default...

### Recursion is quick

To my understanding, the GO compiler is a pretty smart beasty. That does not fully explain why this is happening though.

The two flavours are recursion and a loop with bitwise evaluation.

My first impression was, building a loop is fast. Let's go... Nope... With the 'real' input, it takes about 5 minutes to run. And this is with the use of goroutines!

So, next step, optimise the crap out of it, which did lead to some minor speed increases, but nothing to write home about. It's also how I noticed the bigger than expected penalty on logging statements.

Well... Just because I found the results lacking, I figured I'd at least make the comparison with recursion, even if just to show that it is just a complex problem, which takes a bit of time.

Lo and behold, of course the code looks cleaner, no surprise there. But it outperforms the looping method by a ridiculous factor.

| Concept        | Logging | ns/op       | factor   |
|----------------|---------|-------------|----------|
| Loop & botwise | yes     | 17634369985 |          |
| Recursion      | yes     | 104103530   | ~169 x   |
| -----          |         |             |          |
| Loop & bitwise | no      | 289249480   |          |
| Recursion      | no      | 1367284     | ~211 x   |
| -----          |         |             |          |
| Loop & bitwise | yes     | 17634369985 |          |
| Recursion      | no      | 1367284     | ~12265 x |

So, why?

No idea. Or at least, not why it's so extreme of a difference.
The looping and bitwise operations are all relatively easy operations for a CPU.
- Slice indexing (about x 5 compared to  the recurions method)
    - This could be optimised to x 4 as I implemented a no-op case to keep the bitwise operations a hair simpler.
- Looping with little code in between (little paging between iterations)
- Low nr of 'code paths'/branching (about 3/4 compared to the recurions method)
- Bitwise operations (..., if a language or CPU can't do that efficiently it should just go home....)

On the flip side, the recursion method had these characteristics:
- Simple to read & reason
- Tail-call optimisable-ish
    - It forks into 3, which should be a mark against
- Stack allocation heavy, big memory footprint
- Many re-slices of input data (only pointer cloning though, not values)

So, what have I learned? That the earth is a cube, grass is purple and the atmosphere is a solid substance...
