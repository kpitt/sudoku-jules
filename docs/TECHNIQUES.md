# Deductive Sudoku Techniques

The following is a list of all deductive Sudoku techniques implemented in the
[Hodoku solver](https://hodoku.sourceforge.net/en/techniques.php), in the default
order that they are checked in the solver loop. Each entry includes the Technique
name, the Difficulty level of the technique, and the Status (default techniques
are checked by default, optional techniques are not checked unless explicitly
enabled). The ID column defines the 4-digit identifiers used in the "Hodoku
Library Format" test case format to identify the target technique.

| Order | ID   | Technique                            | Difficulty   | Status   |
| ----: | :--- | :----------------------------------- | :----------- | :------- |
| 1     | 0000 | Full House                           | Easy         | default  |
| 2     | 0003 | Naked Single                         | Easy         | default  |
| 3     | 0002 | Hidden Single                        | Easy         | default  |
| 4     | 0110 | Locked Pair                          | Medium       | default  |
| 5     | 0111 | Locked Triple                        | Medium       | default  |
| 6     | 0100 | Locked Candidates Type 1 (Pointing)  | Medium       | default  |
| 7     | 0101 | Locked Candidates Type 2 (Claiming)  | Medium       | default  |
| 8     | 0200 | Naked Pair                           | Medium       | default  |
| 9     | 0201 | Naked Triple                         | Medium       | default  |
| 10    | 0210 | Hidden Pair                          | Medium       | default  |
| 11    | 0211 | Hidden Triple                        | Medium       | default  |
| 12    | 0202 | Naked Quadruple                      | Hard         | default  |
| 13    | 0212 | Hidden Quadruple                     | Hard         | default  |
| 14    | 0300 | X-Wing                               | Hard         | default  |
| 15    | 0301 | Swordfish                            | Hard         | default  |
| 16    | 0302 | Jellyfish                            | Hard         | default  |
| 17    | 0303 | Squirmbag                            | Unfair       | optional |
| 18    | 0304 | Whale                                | Unfair       | optional |
| 19    | 0305 | Leviathan                            | Unfair       | optional |
| 20    | 0703 | Remote Pair                          | Hard         | default  |
| 21    | 0610 | BUG + 1                              | Hard         | default  |
| 22    | 0400 | Skyscraper                           | Hard         | default  |
| 23    | 0401 | Two String Kite                      | Hard         | default  |
| 24    | 0403 | Turbot Fish                          | Hard         | default  |
| 25    | 0402 | Empty Rectangle                      | Hard         | default  |
| 26    | 0803 | W-Wing                               | Hard         | default  |
| 27    | 0800 | XY-Wing                              | Hard         | default  |
| 28    | 0801 | XYZ-Wing                             | Hard         | default  |
| 29    | 0600 | Uniqueness Test 1                    | Hard         | default  |
| 30    | 0601 | Uniqueness Test 2                    | Hard         | default  |
| 31    | 0602 | Uniqueness Test 3                    | Hard         | default  |
| 32    | 0603 | Uniqueness Test 4                    | Hard         | default  |
| 33    | 0604 | Uniqueness Test 5                    | Hard         | default  |
| 34    | 0605 | Uniqueness Test 6                    | Hard         | default  |
| 35    | 0606 | Hidden Rectangle                     | Hard         | default  |
| 36    | 0607 | Avoidable Rectangle Type 1           | Hard         | default  |
| 37    | 0608 | Avoidable Rectangle Type 2           | Hard         | default  |
| 38    | 0310 | Finned X-Wing                        | Hard         | default  |
| 39    | 0320 | Sashimi X-Wing                       | Hard         | default  |
| 40    | 0311 | Finned Swordfish                     | Unfair       | default  |
| 41    | 0321 | Sashimi Swordfish                    | Unfair       | default  |
| 42    | 0312 | Finned Jellyfish                     | Unfair       | default  |
| 43    | 0322 | Sashimi Jellyfish                    | Unfair       | default  |
| 44    | 0313 | Finned Squirmbag                     | Unfair       | optional |
| 45    | 0323 | Sashimi Squirmbag                    | Unfair       | optional |
| 46    | 0314 | Finned Whale                         | Unfair       | optional |
| 47    | 0324 | Sashimi Whale                        | Unfair       | optional |
| 48    | 0315 | Finned Leviathan                     | Unfair       | optional |
| 49    | 0325 | Sashimi Leviathan                    | Unfair       | optional |
| 50    | 1101 | Sue de Coq                           | Unfair       | default  |
| 51    | 0500 | Simple Colors                        | Hard         | default  |
| 52    | 0502 | Multi Colors                         | Hard         | default  |
| 53    | 0701 | X-Chain                              | Unfair       | default  |
| 54    | 0702 | XY-Chain                             | Unfair       | default  |
| 55    | 0706 | Nice Loop                            | Unfair       | default  |
| 56    | 0709 | Grouped Nice Loop                    | Unfair       | default  |
| 57    | 0901 | ALS-XZ                               | Unfair       | default  |
| 58    | 0902 | ALS-XY-Wing                          | Unfair       | default  |
| 59    | 0903 | ALS-XY-Chain                         | Unfair       | default  |
| 60    | 0904 | Death Blossom                        | Unfair       | optional |
| 61    | 0330 | Franken X-Wing                       | Unfair       | default  |
| 62    | 0331 | Franken Swordfish                    | Unfair       | default  |
| 63    | 0332 | Franken Jellyfish                    | Unfair       | optional |
| 64    | 0333 | Franken Squirmbag                    | Extreme      | optional |
| 65    | 0334 | Franken Whale                        | Extreme      | optional |
| 66    | 0335 | Franken Leviathan                    | Extreme      | optional |
| 67    | 0340 | Finned Franken X-Wing                | Unfair       | default  |
| 68    | 0341 | Finned Franken Swordfish             | Unfair       | default  |
| 69    | 0342 | Finned Franken Jellyfish             | Unfair       | optional |
| 70    | 0343 | Finned Franken Squirmbag             | Extreme      | optional |
| 71    | 0344 | Finned Franken Whale                 | Extreme      | optional |
| 72    | 0345 | Finned Franken Leviathan             | Extreme      | optional |
| 73    | 0350 | Mutant X-Wing                        | Extreme      | optional |
| 74    | 0351 | Mutant Swordfish                     | Extreme      | optional |
| 75    | 0352 | Mutant Jellyfish                     | Extreme      | optional |
| 76    | 0353 | Mutant Squirmbag                     | Extreme      | optional |
| 77    | 0354 | Mutant Whale                         | Extreme      | optional |
| 78    | 0355 | Mutant Leviathan                     | Extreme      | optional |
| 79    | 0360 | Finned Mutant X-Wing                 | Extreme      | optional |
| 80    | 0361 | Finned Mutant Swordfish              | Extreme      | optional |
| 81    | 0362 | Finned Mutant Jellyfish              | Extreme      | optional |
| 82    | 0363 | Finned Mutant Squirmbag              | Extreme      | optional |
| 83    | 0364 | Finned Mutant Whale                  | Extreme      | optional |
| 84    | 0364 | Finned Mutant Leviathan              | Extreme      | optional |
| 85    | 0371 | Kraken Fish                          | Extreme      | optional |
| 86    | 1301 | Forcing Chain                        | Extreme      | default  |
| 87    | 1303 | Forcing Net                          | Extreme      | default  |
| 88    | 1201 | Template Set                         | Extreme      | optional |
| 89    | 1202 | Template Delete                      | Extreme      | optional |
| 90    | xxxx | Brute Force                          | Extreme      | default  |
| 91    | xxxx | Incomplete                           | Incomplete   | optional |
| 92    | xxxx | Give Up                              | Extreme      | default  |
