# ida-subdumper

**[Deprecated]** Please refer to [[Get pseudo code/project from binary with DWARF debug info]](#get-pseudo-codeproject-from-binary-with-dwarf-debug-info) for a better solution.

Dump subroutine pseudocodes from LST file produced by IDA Pro.

1. Open IDA Pro and load the binary file, choose import file names/line numbers (DWARF info found -> Import file names/line numbers), and analyze the binary.
2. Wait for the analysis to finish, and generate the LST file (File -> Produce file -> Create LST file...).
3. Run `go run . {target}` with the target binary file (the path to the LST file is `{target}.lst`), and `{target}.sources.json` will be generated.
4. Edit script `subdumper.py` (replace `{target}` with the target binary file, and define your custom filter function), and load it into IDA Pro (File -> Script file...), and the pseudocodes grouped by source files will be dumped into the `{target}.sources/` directory.

## Get pseudo code/project from binary with DWARF debug info

### Requirements

- A binary with DWARF debug information.
- A tool or library for parsing DWARF debug information. (_e.g._ [go-dwarf](https://github.com/blacktop/go-dwarf))
- A tool or library for disassembling binary and generating pseudocode. (_e.g._ [IDA Pro](https://www.hex-rays.com/ida-pro/), [Ghidra](https://ghidra-sre.org/), [retdec](https://retdec.com/))

### Steps by steps

1. Parse DWARF debug info and get line entries with the following information: address, file, line, column, discriminator. (_e.g._ [line_test.go](https://github.com/blacktop/go-dwarf/blob/main/line_test.go))
2. Group line entries by file and sort them by address.
3. Load binary into the disassembler (no need to analyze or produce LST file).
4. Parse addresses and locate function boundaries.
5. Generate pseudocodes for each function and save them into files in line order.
6. Don't forget to produce the C headers file for structures.
