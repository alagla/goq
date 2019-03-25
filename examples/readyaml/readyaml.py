# example how to parse YAML file

import yaml
from datetime import datetime, timezone

fname = "QuplaTests.yml"

def main():
    print("reading file " + fname)
    fin = open(fname, "r")
    start = unixnow()
    yaml_string = fin.read()
    module_yaml_ir = yaml.load(yaml_string)

    print("module '{}' loaded successfully in {} seconds".format(module_yaml_ir.get("module"), unixnow()-start))
    print("  global type definitions: {}".format(len(module_yaml_ir.get("types"))))
    print("  LUTs: {}".format(len(module_yaml_ir.get("luts"))))
    print("  functions: {}".format(len(module_yaml_ir.get("functions"))))
    print("  executable statements: {}".format(len(module_yaml_ir.get("execs"))))

def unixnow():
    return int(datetime.now(tz=timezone.utc).timestamp())

if __name__ == '__main__':
    main()