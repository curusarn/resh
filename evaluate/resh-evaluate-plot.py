#!/usr/bin/env python3

import sys
import json
from collections import defaultdict
import matplotlib.pyplot as plt
import matplotlib.path as mpath
import numpy as np


def addRank(data):
    return list(enumerate(data, start=1))


data = json.load(sys.stdin)
# for strategy in data["Strategies"]:
#     print(json.dumps(strategy))

cmd_count = defaultdict(int)
cmdLine_count = defaultdict(int)

for record in data["Records"]:
    cmd_count[record["firstWord"]] += 1
    cmdLine_count[record["cmdLine"]] += 1

cmdTmp = sorted(cmd_count.items(), key=lambda x: x[1], reverse=True)[:50]
cmdFrq = list(map(lambda x: x[1] / cmdTmp[0][1], cmdTmp))

cmdLineTmp = sorted(cmdLine_count.items(), key=lambda x: x[1], reverse=True)[:50]
cmdLineFrq = list(map(lambda x: x[1] / cmdLineTmp[0][1], cmdLineTmp))

print(cmdFrq)
print("#################")
#print(cmdLineFrq_rank)

plt.plot(range(1, len(cmdFrq)+1), cmdFrq, 'o-')
plt.title("Command frequency")
plt.yticks()
#plt.xticks(range(1, len(cmdFrq)+1))
plt.show()

plt.plot(range(1, len(cmdLineFrq)+1), cmdLineFrq, 'o-')
plt.title("Commandline frequency")
plt.show()