#!/usr/bin/env python3

import sys
import json
from collections import defaultdict
import matplotlib.pyplot as plt
import matplotlib.path as mpath
import numpy as np
from graphviz import Digraph

PLOT_WIDTH = 10 # inches
PLOT_HEIGHT = 7 # inches

PLOT_SIZE_zipf = 20

data = json.load(sys.stdin)
# for strategy in data["Strategies"]:
#     print(json.dumps(strategy))


def zipf(length):
    return list(map(lambda x: 1/2**x, range(0, length)))


def trim(text, length, add_elipse=True):
    if add_elipse and len(text) > length:
        return text[:length-1] + "…"
    return text[:length]


# Figure 3.1. The normalized command frequency, compared with Zipf.
def plot_cmdLineFrq_rank(plotSize=PLOT_SIZE_zipf, show_labels=False):
    cmdLine_count = defaultdict(int)
    for record in data["Records"]:
        if record["invalid"]:
            continue

        cmdLine_count[record["cmdLine"]] += 1

    tmp = sorted(cmdLine_count.items(), key=lambda x: x[1], reverse=True)[:plotSize]
    cmdLineFrq = list(map(lambda x: x[1] / tmp[0][1], tmp))
    labels = list(map(lambda x: trim(x[0], 7), tmp))

    ranks = range(1, len(cmdLineFrq)+1)
    plt.figure(figsize=(PLOT_WIDTH, PLOT_HEIGHT))
    plt.plot(ranks, zipf(len(ranks)), '-')
    plt.plot(ranks, cmdLineFrq, 'o-')
    plt.title("Commandline frequency / rank")
    plt.ylabel("Normalized commandline frequency")
    plt.xlabel("Commandline rank")
    plt.legend(("Zipf", "Commandline"), loc="best")
    if show_labels:
        plt.xticks(ranks, labels, rotation=-60)
    # TODO: make xticks integral
    plt.show()


# similar to ~ Figure 3.1. The normalized command frequency, compared with Zipf.
def plot_cmdFrq_rank(plotSize=PLOT_SIZE_zipf, show_labels=False):
    cmd_count = defaultdict(int)
    for record in data["Records"]:
        if record["invalid"]:
            continue

        cmd = record["firstWord"]
        if cmd == "":
            continue
        cmd_count[cmd] += 1

    tmp = sorted(cmd_count.items(), key=lambda x: x[1], reverse=True)[:plotSize]
    cmdFrq = list(map(lambda x: x[1] / tmp[0][1], tmp))
    labels = list(map(lambda x: trim(x[0], 7), tmp))

    ranks = range(1, len(cmdFrq)+1)
    plt.figure(figsize=(PLOT_WIDTH, PLOT_HEIGHT))
    plt.plot(ranks, zipf(len(ranks)), 'o-')
    plt.plot(ranks, cmdFrq, 'o-')
    plt.title("Command frequency / rank")
    plt.ylabel("Normalized command frequency")
    plt.xlabel("Command rank")
    plt.legend(("Zipf", "Command"), loc="best")
    if show_labels:
        plt.xticks(ranks, labels, rotation=-60)
    # TODO: make xticks integral
    plt.show()

# Figure 3.2. Command vocabulary size vs. the number of command lines entered for four individuals.
def plot_cmdVocabularySize_cmdLinesEntered():
    cmd_vocabulary = set()
    y_cmd_count = [0]
    for record in data["Records"]:
        if record["invalid"]:
            continue

        cmd = record["firstWord"]
        if cmd in cmd_vocabulary:
            # repeat last value
            y_cmd_count.append(y_cmd_count[-1])
        else:
            cmd_vocabulary.add(cmd)  
            # append last value +1
            y_cmd_count.append(y_cmd_count[-1] + 1)

    print(cmd_vocabulary)
    x_cmds_entered = range(0, len(y_cmd_count))

    plt.figure(figsize=(PLOT_WIDTH, PLOT_HEIGHT))
    plt.plot(x_cmds_entered, y_cmd_count, '-')
    plt.title("Command vocabulary size vs. the number of command lines entered")
    plt.ylabel("Command vocabulary size")
    plt.xlabel("# of command lines entered")
    plt.show()

# Figure 3.3. Sequential structure of UNIX command usage, from Figure 4 in Hanson et al. (1984).
#       Ball diameters are proportional to stationary probability. Lines indicate significant dependencies,
#       solid ones being more probable (p < .0001) and dashed ones less probable (.005 < p < .0001).
def graphviz_cmdSequences(cmd_displayTreshold=20, edge_displayTreshold=0.03):
    cmd_count = defaultdict(int)
    cmdSeq_count = defaultdict(lambda: defaultdict(int))
    cmd_id = dict()
    prev_cmd = "_SESSION_INIT_" # XXX: not actually session init yet
    cmd_id[prev_cmd] = str(-1) 
    for x, record in enumerate(data["Records"]):
        if record["invalid"]:
            continue

        cmd = record["firstWord"]
        cmdSeq_count[prev_cmd][cmd] += 1
        cmd_count[cmd] += 1
        cmd_id[cmd] = str(x)
        prev_cmd = cmd

    dot = Digraph(comment="Command sequences", graph_attr={'overlap':'scale', 'splines':'true'})

    # for cmd_entry in cmdSeq_count.items():
    #     cmd, seq = cmd_entry

    #     if cmd_count[cmd] < cmd_displayTreshold:
    #         continue
    #     
    #     dot.node(cmd_id[cmd], cmd)

    for cmd_entry in cmdSeq_count.items():
        cmd, seq = cmd_entry

        count = cmd_count[cmd]
        if count < cmd_displayTreshold:
            continue

        for seq_entry in seq.items():
            cmd2, seq_count = seq_entry
            relative_seq_count = seq_count / count

            if cmd_count[cmd2] < cmd_displayTreshold:
                continue
            if relative_seq_count < edge_displayTreshold:
                continue
            
            for id_, cmd_ in ((cmd_id[cmd], cmd), (cmd_id[cmd2], cmd2)):
                count_ = cmd_count[cmd_]
                scale_ = count_ / (cmd_displayTreshold)
                width_ = str(0.08*scale_) 
                fontsize_ = str(1*scale_)
                if scale_ < 12:
                    dot.node(id_, '', shape='circle', fixedsize='true', fontname='bold',
                            width=width_, fontsize='12', forcelabels='true', xlabel=cmd_)
                else:
                    dot.node(id_, cmd_, shape='circle', fixedsize='true', fontname='bold',
                            width=width_, fontsize=fontsize_, forcelabels='true')

            
            # 1.0 is max
            scale_ = seq_count / cmd_count[cmd]
            penwidth_ = str(0.5 + 4.5 * scale_)
            #penwidth_bold_ = str(8 * scale_)
            if scale_ > 0.5:
                dot.edge(cmd_id[cmd], cmd_id[cmd2], constraint='false', splines='curved',
                         penwidth=penwidth_, style='bold')
            elif scale_ > 0.2:
                dot.edge(cmd_id[cmd], cmd_id[cmd2], constraint='false', splines='curved',
                         penwidth=penwidth_, arrowhead='open')
            elif scale_ > 0.1:
                dot.edge(cmd_id[cmd], cmd_id[cmd2], constraint='false', splines='curved',
                         penwidth=penwidth_, style='dashed', arrowhead='open')
            else:
                dot.edge(cmd_id[cmd], cmd_id[cmd2], constraint='false', splines='curved',
                         penwidth=penwidth_, style='dotted', arrowhead='empty')

    dot.render('/tmp/resh-graphviz-cmdSeq.gv', view=False)

graphviz_cmdSequences()
# plot_cmdVocabularySize_cmdLinesEntered()
# plot_cmdLineFrq_rank()
# plot_cmdFrq_rank()


# be careful and check if labels fit the display