#!/usr/bin/env python3


import traceback
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

DATA_records = []
DATA_records_by_session = defaultdict(list) 
for user in data["UsersRecords"]:
    for device in user["Devices"]:
        for record in device["Records"]:
            if record["invalid"]:
                continue
            
            DATA_records.append(record)
            DATA_records_by_session[record["sessionPid"]].append(record)

DATA_records = list(sorted(DATA_records, key=lambda x: x["realtimeBeforeLocal"]))

for pid, session in DATA_records_by_session.items():
    session = list(sorted(session, key=lambda x: x["realtimeBeforeLocal"]))


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
    for record in DATA_records:
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
    for record in DATA_records:
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
    for record in DATA_records:
        cmd = record["firstWord"]
        if cmd in cmd_vocabulary:
            # repeat last value
            y_cmd_count.append(y_cmd_count[-1])
        else:
            cmd_vocabulary.add(cmd)  
            # append last value +1
            y_cmd_count.append(y_cmd_count[-1] + 1)

    # print(cmd_vocabulary)
    x_cmds_entered = range(0, len(y_cmd_count))

    plt.figure(figsize=(PLOT_WIDTH, PLOT_HEIGHT))
    plt.plot(x_cmds_entered, y_cmd_count, '-')
    plt.title("Command vocabulary size vs. the number of command lines entered")
    plt.ylabel("Command vocabulary size")
    plt.xlabel("# of command lines entered")
    plt.show()

# Figure 5.6. Command line vocabulary size vs. the number of commands entered for four typical individuals.
def plot_cmdLineVocabularySize_cmdLinesEntered():
    cmdLine_vocabulary = set()
    y_cmdLine_count = [0]
    for record in DATA_records:
        cmdLine = record["cmdLine"]
        if cmdLine in cmdLine_vocabulary:
            # repeat last value
            y_cmdLine_count.append(y_cmdLine_count[-1])
        else:
            cmdLine_vocabulary.add(cmdLine)  
            # append last value +1
            y_cmdLine_count.append(y_cmdLine_count[-1] + 1)

    # print(cmdLine_vocabulary)
    x_cmdLines_entered = range(0, len(y_cmdLine_count))

    plt.figure(figsize=(PLOT_WIDTH, PLOT_HEIGHT))
    plt.plot(x_cmdLines_entered, y_cmdLine_count, '-')
    plt.title("Command line vocabulary size vs. the number of command lines entered")
    plt.ylabel("Command line vocabulary size")
    plt.xlabel("# of command lines entered")
    plt.show()

# Figure 3.3. Sequential structure of UNIX command usage, from Figure 4 in Hanson et al. (1984).
#       Ball diameters are proportional to stationary probability. Lines indicate significant dependencies,
#       solid ones being more probable (p < .0001) and dashed ones less probable (.005 < p < .0001).
def graph_cmdSequences(node_count=33, edge_minValue=0.05):
    START_CMD = "_start_"
    cmd_count = defaultdict(int)
    cmdSeq_count = defaultdict(lambda: defaultdict(int))
    cmd_id = dict()
    x = 0
    cmd_id[START_CMD] = str(x) 
    for pid, session in DATA_records_by_session.items():
        cmd_count[START_CMD] += 1
        prev_cmd = START_CMD
        for record in session:
            cmd = record["firstWord"]
            cmdSeq_count[prev_cmd][cmd] += 1
            cmd_count[cmd] += 1
            if cmd not in cmd_id:
                x += 1
                cmd_id[cmd] = str(x)
            prev_cmd = cmd

    # get `node_count` of largest nodes
    sorted_cmd_count = sorted(cmd_count.items(), key=lambda x: x[1], reverse=True)
    print(sorted_cmd_count)
    cmds_to_graph = list(map(lambda x: x[0], sorted_cmd_count))[:node_count]

    # use 3 biggest nodes as a reference point for scaling
    biggest_node = cmd_count[cmds_to_graph[0]]
    nd_biggest_node = cmd_count[cmds_to_graph[1]]
    rd_biggest_node = cmd_count[cmds_to_graph[1]]
    count2scale_coef = 3 / (biggest_node + nd_biggest_node + rd_biggest_node)

    # scaling constant
    #       affects node size and node label
    base_scaling_factor = 21
    # extra scaling for experiments - not really useful imho
    #       affects everything nodes, edges, node labels, treshold for turning label into xlabel, xlabel size, ...
    extra_scaling_factor = 1.0 
    for x in range(0, 10):
        # graphviz is not the most reliable piece of software
        #       -> retry on fail but scale nodes down by 1%
        scaling_factor = base_scaling_factor * (1 - x * 0.01)

        # overlap: scale -> solve overlap by scaling the graph
        # overlap_shrink -> try to shrink the graph a bit after you are done
        # splines -> don't draw edges over nodes
        # sep: 2.5 -> assume that nodes are 2.5 inches larger
        graph_attr={'overlap':'scale', 'overlap_shrink':'true',
                    'splines':'true', 'sep':'0.25'}
        graph = Digraph(name='command_sequentiality', engine='neato', graph_attr=graph_attr)

        # iterate over all nodes
        for cmd in cmds_to_graph:
            seq = cmdSeq_count[cmd]
            count = cmd_count[cmd]

            # iterate over all "following" commands (for each node)
            for seq_entry in seq.items():
                cmd2, seq_count = seq_entry
                relative_seq_count = seq_count / count

                # check if "follow" command is supposed to be in the graph
                if cmd2 not in cmds_to_graph:
                    continue
                # check if the edge value is high enough
                if relative_seq_count < edge_minValue:
                    continue
                
                # create starting node and end node for the edge
                #       duplicates don't matter 
                for id_, cmd_ in ((cmd_id[cmd], cmd), (cmd_id[cmd2], cmd2)):
                    count_ = cmd_count[cmd_]
                    scale_ = count_ * count2scale_coef * scaling_factor * extra_scaling_factor
                    width_ = 0.08 * scale_
                    fontsize_ = 8.5 * scale_ / (len(cmd_) + 3)

                    width_ = str(width_) 
                    if fontsize_ < 12 * extra_scaling_factor:
                        graph.node(id_, ' ', shape='circle', fixedsize='true', fontname='monospace bold',
                                width=width_, fontsize=str(12 * extra_scaling_factor), forcelabels='true', xlabel=cmd_)
                    else:
                        fontsize_ = str(fontsize_)
                        graph.node(id_, cmd_, shape='circle', fixedsize='true', fontname='monospace bold',
                                width=width_, fontsize=fontsize_, forcelabels='true', labelloc='c')
                
                # value of the edge (percentage) 1.0 is max
                scale_ = seq_count / cmd_count[cmd]
                penwidth_ = str((0.5 + 4.5 * scale_) * extra_scaling_factor)
                #penwidth_bold_ = str(8 * scale_)
                if scale_ > 0.5:
                    graph.edge(cmd_id[cmd], cmd_id[cmd2], constraint='true', splines='curved',
                            penwidth=penwidth_, style='bold')
                elif scale_ > 0.2:
                    graph.edge(cmd_id[cmd], cmd_id[cmd2], constraint='true', splines='curved',
                            penwidth=penwidth_, arrowhead='open')
                elif scale_ > 0.1:
                    graph.edge(cmd_id[cmd], cmd_id[cmd2], constraint='true', splines='curved',
                            penwidth=penwidth_, style='dashed', arrowhead='open')
                else:
                    graph.edge(cmd_id[cmd], cmd_id[cmd2], constraint='false', splines='curved',
                            penwidth=penwidth_, style='dotted', arrowhead='empty')

        # graphviz sometimes fails - see above
        try:
            graph.view()
            # graph.render('/tmp/resh-graphviz-cmdSeq.gv', view=True)
            break
        except Exception as e:
            trace = traceback.format_exc()
            print("GRAPHVIZ EXCEPTION: <{}>\nGRAPHVIZ TRACE: <{}>".format(str(e), trace))


def plot_strategies_matches(plot_size=50, selected_strategies=[]):
    plt.figure(figsize=(PLOT_WIDTH, PLOT_HEIGHT))
    plt.title("Matches at distance")
    plt.ylabel('%' + " of matches")
    plt.xlabel("Distance")
    legend = []
    for strategy in data["Strategies"]:
        strategy_title = strategy["Title"]
        # strategy_description = strategy["Description"]

        if len(selected_strategies) and strategy_title not in selected_strategies:
            continue

        dataPoint_count = 0
        matches = [0] * plot_size
        matches_total = 0
        charsRecalled = [0] * plot_size
        charsRecalled_total = 0
        
        for match in strategy["Matches"]:
            dataPoint_count += 1

            if not match["Match"]:
                continue

            chars = match["CharsRecalled"]
            charsRecalled_total += chars 
            matches_total += 1

            dist = match["Distance"]  
            if dist > plot_size:
                continue

            matches[dist-1] += 1
            charsRecalled[dist-1] += chars
            

        acc = 0
        matches_cumulative = []
        for x in matches:
            acc += x
            matches_cumulative.append(acc)
        matches_cumulative.append(matches_total)
        matches_percent = list(map(lambda x: 100 * x / dataPoint_count, matches_cumulative))

        x_values = range(1, plot_size+2)
        plt.plot(x_values, matches_percent, 'o-')
        legend.append(strategy_title)


    x_ticks = list(range(1, plot_size+1, 2))
    x_labels = x_ticks[:]
    x_ticks.append(plot_size+1)
    x_labels.append("total")
    plt.xticks(x_ticks, x_labels)
    plt.legend(legend, loc="best")
    plt.show()



def plot_strategies_charsRecalled(plot_size=50, selected_strategies=[]):
    plt.figure(figsize=(PLOT_WIDTH, PLOT_HEIGHT))
    plt.title("Average characters recalled at distance")
    plt.ylabel("Average characters recalled")
    plt.xlabel("Distance")
    legend = []
    for strategy in data["Strategies"]:
        strategy_title = strategy["Title"]
        # strategy_description = strategy["Description"]

        if len(selected_strategies) and strategy_title not in selected_strategies:
            continue

        dataPoint_count = 0
        matches = [0] * plot_size
        matches_total = 0
        charsRecalled = [0] * plot_size
        charsRecalled_total = 0
        
        for match in strategy["Matches"]:
            dataPoint_count += 1

            if not match["Match"]:
                continue

            chars = match["CharsRecalled"]
            charsRecalled_total += chars 
            matches_total += 1

            dist = match["Distance"]  
            if dist > plot_size:
                continue

            matches[dist-1] += 1
            charsRecalled[dist-1] += chars
            

        acc = 0
        charsRecalled_cumulative = []
        for x in charsRecalled:
            acc += x
            charsRecalled_cumulative.append(acc)
        charsRecalled_cumulative.append(charsRecalled_total)
        charsRecalled_average = list(map(lambda x: x / dataPoint_count, charsRecalled_cumulative))

        x_values = range(1, plot_size+2)
        plt.plot(x_values, charsRecalled_average, 'o-')
        legend.append(strategy_title)


    x_ticks = list(range(1, plot_size+1, 2))
    x_labels = x_ticks[:]
    x_ticks.append(plot_size+1)
    x_labels.append("total")
    plt.xticks(x_ticks, x_labels)
    plt.legend(legend, loc="best")
    plt.show()


        
graph_cmdSequences()
graph_cmdSequences(node_count=28, edge_minValue=0.06)

plot_cmdLineFrq_rank()
# plot_cmdFrq_rank()
        
plot_cmdLineVocabularySize_cmdLinesEntered()
# plot_cmdVocabularySize_cmdLinesEntered()

plot_strategies_matches()
plot_strategies_charsRecalled()

# be careful and check if labels fit the display