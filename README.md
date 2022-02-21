# mkchop
Split crypto phrase into parts and combine back using Shamir split algorithm.
Learn more about it [here](https://youtu.be/K54ildEW9-Q)

## disclaimer
>The use of this tool does not guarantee security or suitability
for any particular use. Please review the code and use at your own risk.

## installation
This step assumes you have [Go compiler toolchain](https://go.dev/dl/)
installed on your system.

Download the code to a folder and cd to the folder, then run
```bash
go install
```
Install shell completion. For instance `bash` completion can be installed
by adding following line to your `.bashrc`:
```bash
source <(mkchop completion bash)
```

## split input into parts
Split input into parts such that a minimum threshold of them can be used
later to reconstruct the original input.

The input is read from STDIN and output is a set of base58 strings as shown
below. The default values of splits is 5 and the default minimum threshold to
reconstruct back is 3.
```bash
mkchop split 
Enter input string: maple pacific mirror blue car country computer store air
VsCUMRD1uLM68a66h3ZKbiTtcnpBF6TpbVC2BddTRbQ9cC37Wtcov5ngEowAgM2YZHyjULg2Am7CTm
H6CeNRJ73jKVaD3LTvKQEeh79U54tXCjajbSxcceLiuDVAEFFtNH9m9VoRVMMTzFcCG6eEQzJxnqEE
Pdf8iwJgGFjavsNoEVZqSb4QU4uwRUpr1KQFxNAQs1uuSr2412XpD6giwWB3st9oDcqeHPFV6oyrPe
Gh4nSSXqn4Sry84wApF7MPAGGRJyicX7cMGqdL2sBqTzF4GKxsno4unKEmNk8RbSNFd99RZT3VPDr3
TqyPmxrdMko3UELuzf8tsBuicGXDtc9dfHAN61LbgNA91j1QM2rxdgn8BR8iKuXCvctUZ3XCZWeQVq
```

## combine parts back to original input
At least a minimum number of parts are required (in any order) to construct
the original input back as shown below:
```bash
mkchop combine 
Enter part 1 of 3: TqyPmxrdMko3UELuzf8tsBuicGXDtc9dfHAN61LbgNA91j1QM2rxdgn8BR8iKuXCvctUZ3XCZWeQVq
Enter part 2 of 3: H6CeNRJ73jKVaD3LTvKQEeh79U54tXCjajbSxcceLiuDVAEFFtNH9m9VoRVMMTzFcCG6eEQzJxnqEE
Enter part 3 of 3: Gh4nSSXqn4Sry84wApF7MPAGGRJyicX7cMGqdL2sBqTzF4GKxsno4unKEmNk8RbSNFd99RZT3VPDr3
maple pacific mirror blue car country computer store air
```
