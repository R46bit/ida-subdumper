#!/usr/bin/env python3

import json
import os

import idaapi


TARGET = "{target}"


def filter(source):
    if source == "":
        return True
    return False


def decompile(addr):
    try:
        cfunc = idaapi.decompile(addr)
    except idaapi.DecompilationFailure:
        cfunc = None
    if cfunc is None:
        return None
    lines = cfunc.get_pseudocode()
    retlines = []
    for lnnum in range(len(lines)):
        retlines.append(idaapi.tag_remove(lines[lnnum].line))
    return '\n'.join(retlines)


with open(TARGET+".sources.json") as json_file:
    sources = json.load(json_file)
    for source, subroutines in sources.items():
        if filter(source):
            continue
        print("dump", source)
        name = TARGET + ".sources" + source
        os.makedirs(os.path.dirname(name), exist_ok=True)
        with open(name, "w") as f:
            f.write(f"// File: {source}\n")
            for subroutine in subroutines:
                f.write("\n")
                f.write(f"// Line {subroutine['line']}: range {subroutine['addr']}\n")
                func = decompile(int("0x"+subroutine['addr'][:16], 16))
                if func is not None:
                    f.write(func+";")
                else:
                    f.write(subroutine['func']+";")
                f.write("\n")
