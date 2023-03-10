# ida-subdumper

Dump subroutine pseudocodes from LST file produced by IDA Pro.

1. Open IDA Pro and load the binary file, choose import file names/line numbers (DWARF info found -> Import file names/line numbers), and analyze the binary.
2. Wait for the analysis to finish, and generate the LST file (File -> Produce file -> Create LST file...).
3. Run `go run . {target}` with the target binary file (the path to the LST file is `{target}.lst`), and `{target}.sources.json` will be generated.
4. Edit script `subdumper.py` (replace `{target}` with the target binary file, and define your custom filter function), and load it into IDA Pro (File -> Script file...), and the pseudocodes grouped by source files will be dumped into the `{target}.sources/` directory.
