#!/usr/bin/env python3


import traceback
import sys
import json
from collections import defaultdict
import numpy as np
from graphviz import Digraph
from datetime import datetime

from matplotlib import rcParams
rcParams['font.family'] = 'serif'
# rcParams['font.serif'] = ['']

import matplotlib.pyplot as plt
import matplotlib.path as mpath
import matplotlib.patches as mpatches

PLOT_WIDTH = 10 # inches
PLOT_HEIGHT = 7 # inches

PLOT_SIZE_zipf = 20

data = json.load(sys.stdin)

DATA_records = []
DATA_records_by_session = defaultdict(list) 
DATA_records_by_user = defaultdict(list) 
for user in data["UsersRecords"]:
    if user["Devices"] is None:
        continue
    for device in user["Devices"]:
        if device["Records"] is None:
            continue
        for record in device["Records"]:
            if "invalid" in record and record["invalid"]:
                continue
            
            DATA_records.append(record)
            DATA_records_by_session[record["seqSessionId"]].append(record)
            DATA_records_by_user[user["Name"] + ":" + device["Name"]].append(record)

DATA_records = list(sorted(DATA_records, key=lambda x: x["realtimeAfterLocal"]))

for pid, session in DATA_records_by_session.items():
    session = list(sorted(session, key=lambda x: x["realtimeAfterLocal"]))

# TODO: this should be a cmdline option
async_draw = True

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
    if async_draw:
        plt.draw()
    else:
        plt.show()


# similar to ~ Figure 3.1. The normalized command frequency, compared with Zipf.
def plot_cmdFrq_rank(plotSize=PLOT_SIZE_zipf, show_labels=False):
    plt.figure(figsize=(PLOT_WIDTH, PLOT_HEIGHT))
    plt.title("Command frequency / rank")
    plt.ylabel("Normalized command frequency")
    plt.xlabel("Command rank")
    legend = []


    cmd_count = defaultdict(int)
    len_records = 0
    for record in DATA_records:
        cmd = record["command"]
        if cmd == "":
            continue
        cmd_count[cmd] += 1
        len_records += 1

    tmp = sorted(cmd_count.items(), key=lambda x: x[1], reverse=True)[:plotSize]
    cmdFrq = list(map(lambda x: x[1] / tmp[0][1], tmp))
    labels = list(map(lambda x: trim(x[0], 7), tmp))

    top100percent = 100 * sum(map(lambda x: x[1], list(cmd_count.items())[:int(1 * len(cmd_count))])) / len_records
    top10percent = 100 * sum(map(lambda x: x[1], list(cmd_count.items())[:int(0.1 * len(cmd_count))])) / len_records
    top20percent = 100 * sum(map(lambda x: x[1], list(cmd_count.items())[:int(0.2 * len(cmd_count))])) / len_records
    print("% ALL: Top {} %% of cmds amounts for {} %% of all command lines".format(100, top100percent))
    print("% ALL: Top {} %% of cmds amounts for {} %% of all command lines".format(10, top10percent))
    print("% ALL: Top {} %% of cmds amounts for {} %% of all command lines".format(20, top20percent))
    ranks = range(1, len(cmdFrq)+1)
    plt.plot(ranks, zipf(len(ranks)), '-')
    legend.append("Zipf distribution")
    plt.plot(ranks, cmdFrq, 'o-')
    legend.append("All subjects")


    for user in DATA_records_by_user.items():
        cmd_count = defaultdict(int)
        len_records = 0
        name, records = user
        for record in records:
            cmd = record["command"]
            if cmd == "":
                continue
            cmd_count[cmd] += 1
            len_records += 1

        tmp = sorted(cmd_count.items(), key=lambda x: x[1], reverse=True)[:plotSize]
        cmdFrq = list(map(lambda x: x[1] / tmp[0][1], tmp))
        labels = list(map(lambda x: trim(x[0], 7), tmp))

        top100percent = 100 * sum(map(lambda x: x[1], list(cmd_count.items())[:int(1 * len(cmd_count))])) / len_records
        top10percent = 100 * sum(map(lambda x: x[1], list(cmd_count.items())[:int(0.1 * len(cmd_count))])) / len_records
        top20percent = 100 * sum(map(lambda x: x[1], list(cmd_count.items())[:int(0.2 * len(cmd_count))])) / len_records
        print("% {}: Top {} %% of cmds amounts for {} %% of all command lines".format(name, 100, top100percent))
        print("% {}: Top {} %% of cmds amounts for {} %% of all command lines".format(name, 10, top10percent))
        print("% {}: Top {} %% of cmds amounts for {} %% of all command lines".format(name, 20, top20percent))
        ranks = range(1, len(cmdFrq)+1)
        plt.plot(ranks, cmdFrq, 'o-')
        legend.append("{} (sanitize!)".format(name))

    plt.legend(legend, loc="best")

    if show_labels:
        plt.xticks(ranks, labels, rotation=-60)
    # TODO: make xticks integral
    if async_draw:
        plt.draw()
    else:
        plt.show()

# Figure 3.2. Command vocabulary size vs. the number of command lines entered for four individuals.
def plot_cmdVocabularySize_cmdLinesEntered():
    plt.figure(figsize=(PLOT_WIDTH, PLOT_HEIGHT))
    plt.title("Command vocabulary size vs. the number of command lines entered")
    plt.ylabel("Command vocabulary size")
    plt.xlabel("# of command lines entered")
    legend = []

    # x_count = max(map(lambda x: len(x[1]), DATA_records_by_user.items()))
    # x_values = range(0, x_count)  
    for user in DATA_records_by_user.items():
        new_cmds_after_1k = 0
        new_cmds_after_2k = 0
        new_cmds_after_3k = 0
        cmd_vocabulary = set()
        y_cmd_count = [0]
        name, records = user
        for record in records:
            cmd = record["command"]
            if cmd == "":
                continue
            if cmd in cmd_vocabulary:
                # repeat last value
                y_cmd_count.append(y_cmd_count[-1])
            else:
                cmd_vocabulary.add(cmd)  
                # append last value +1
                y_cmd_count.append(y_cmd_count[-1] + 1)
                if len(y_cmd_count) > 1000:
                    new_cmds_after_1k+=1
                if len(y_cmd_count) > 2000:
                    new_cmds_after_2k+=1
                if len(y_cmd_count) > 3000:
                    new_cmds_after_3k+=1
        
            if len(y_cmd_count) == 1000:
                print("% {}: Cmd adoption rate at 1k (between 0 and 1k) cmdlines = {}".format(name ,len(cmd_vocabulary) / (len(y_cmd_count))))
            if len(y_cmd_count) == 2000:
                print("% {}: Cmd adoption rate at 2k cmdlines = {}".format(name ,len(cmd_vocabulary) / (len(y_cmd_count))))
                print("% {}: Cmd adoption rate between 1k and 2k cmdlines = {}".format(name ,new_cmds_after_1k / (len(y_cmd_count) - 1000)))
            if len(y_cmd_count) == 3000:
                print("% {}: Cmd adoption rate between 2k and 3k cmdlines = {}".format(name ,new_cmds_after_2k / (len(y_cmd_count) - 2000)))

        print("% {}: New cmd adoption rate after 1k cmdlines = {}".format(name ,new_cmds_after_1k / (len(y_cmd_count) - 1000)))
        print("% {}: New cmd adoption rate after 2k cmdlines = {}".format(name ,new_cmds_after_2k / (len(y_cmd_count) - 2000)))
        print("% {}: New cmd adoption rate after 3k cmdlines = {}".format(name ,new_cmds_after_3k / (len(y_cmd_count) - 3000)))
        x_cmds_entered = range(0, len(y_cmd_count))
        plt.plot(x_cmds_entered, y_cmd_count, '-')
        legend.append(name + " (TODO: sanitize!)")

    # print(cmd_vocabulary)

    plt.legend(legend, loc="best")

    if async_draw:
        plt.draw()
    else:
        plt.show()


def plot_cmdVocabularySize_daily():
    SECONDS_IN_A_DAY = 86400
    plt.figure(figsize=(PLOT_WIDTH, PLOT_HEIGHT))
    plt.title("Command vocabulary size in days")
    plt.ylabel("Command vocabulary size")
    plt.xlabel("Days")
    legend = []

    # x_count = max(map(lambda x: len(x[1]), DATA_records_by_user.items()))
    # x_values = range(0, x_count)  
    for user in DATA_records_by_user.items():
        new_cmds_after_100 = 0
        new_cmds_after_200 = 0
        new_cmds_after_300 = 0
        cmd_vocabulary = set()
        y_cmd_count = [0]
        name, records = user

        cmd_fail_count = 0

        if not len(records):
            print("ERROR: no records for user {}".format(name))
            continue

        first_day = records[0]["realtimeAfter"]
        this_day = first_day

        for record in records:
            cmd = record["command"]
            timestamp = record["realtimeAfter"]

            if cmd == "":
                cmd_fail_count += 1
                continue

            if timestamp >= this_day + SECONDS_IN_A_DAY:
                this_day += SECONDS_IN_A_DAY
                while timestamp >= this_day + SECONDS_IN_A_DAY:
                    y_cmd_count.append(-10)
                    this_day += SECONDS_IN_A_DAY

                y_cmd_count.append(len(cmd_vocabulary))
                cmd_vocabulary = set() # wipes the vocabulary each day

                if len(y_cmd_count) > 100:
                    new_cmds_after_100+=1
                if len(y_cmd_count) > 200:
                    new_cmds_after_200+=1
                if len(y_cmd_count) > 300:
                    new_cmds_after_300+=1

                if len(y_cmd_count) == 100:
                    print("% {}: Cmd adoption rate at 100 days (between 0 and 100 days) = {}".format(name, len(cmd_vocabulary) / (len(y_cmd_count))))
                if len(y_cmd_count) == 200:
                    print("% {}: Cmd adoption rate at 200 days days = {}".format(name, len(cmd_vocabulary) / (len(y_cmd_count))))
                    print("% {}: Cmd adoption rate between 100 and 200 days = {}".format(name, new_cmds_after_100 / (len(y_cmd_count) - 100)))
                if len(y_cmd_count) == 300:
                    print("% {}: Cmd adoption rate between 200 and 300 days = {}".format(name, new_cmds_after_200 / (len(y_cmd_count) - 200)))

            if cmd not in cmd_vocabulary:
                cmd_vocabulary.add(cmd)  
        

        print("% {}: New cmd adoption rate after 100 days = {}".format(name, new_cmds_after_100 / (len(y_cmd_count) - 100)))
        print("% {}: New cmd adoption rate after 200 days = {}".format(name, new_cmds_after_200 / (len(y_cmd_count) - 200)))
        print("% {}: New cmd adoption rate after 300 days = {}".format(name, new_cmds_after_300 / (len(y_cmd_count) - 300)))
        print("% {}: cmd_fail_count = {}".format(name, cmd_fail_count))
        x_cmds_entered = range(0, len(y_cmd_count))
        plt.plot(x_cmds_entered, y_cmd_count, 'o', markersize=2)
        legend.append(name + " (TODO: sanitize!)")

    # print(cmd_vocabulary)

    plt.legend(legend, loc="best")
    plt.ylim(bottom=-5)

    if async_draw:
        plt.draw()
    else:
        plt.show()


def matplotlib_escape(ss):
    ss = ss.replace('$', '\\$')
    return ss


def plot_cmdUsage_in_time(sort_cmds=False, num_cmds=None):
    SECONDS_IN_A_DAY = 86400
    tab_colors = ("tab:blue", "tab:orange", "tab:green", "tab:red", "tab:purple", "tab:brown", "tab:pink", "tab:gray")
    plt.figure(figsize=(PLOT_WIDTH, PLOT_HEIGHT))
    plt.title("Command use in time")
    plt.ylabel("Commands")
    plt.xlabel("Days")
    legend_patches = []

    cmd_ids = {}
    y_labels = []

    all_x_values = []
    all_y_values = []
    all_s_values = [] # size
    all_c_values = [] # color

    x_values = []
    y_values = []
    s_values = [] # size
    c_values = [] # color

    if sort_cmds:
        cmd_count = defaultdict(int)
        for user in DATA_records_by_user.items():
            name, records = user
            for record in records:
                cmd = record["command"]
                cmd_count[cmd] += 1

        sorted_cmds = map(lambda x: x[0], sorted(cmd_count.items(), key=lambda x: x[1], reverse=True))

        for cmd in sorted_cmds:
            cmd_ids[cmd] = len(cmd_ids)
            y_labels.append(matplotlib_escape(cmd))

    
    for user_idx, user in enumerate(DATA_records_by_user.items()):
        name, records = user

        if not len(records):
            print("ERROR: no records for user {}".format(name))
            continue


        first_day = records[0]["realtimeAfter"]
        this_day = first_day
        day_no = 0 
        today_cmds = defaultdict(int) 

        for record in records:
            cmd = record["command"]
            timestamp = record["realtimeAfter"]

            if cmd == "":
                print("NOTICE: Empty cmd for {}".format(record["cmdLine"]))
                continue

            if timestamp >= this_day + SECONDS_IN_A_DAY:
                for item in today_cmds.items():
                    cmd, count = item
                    cmd_id = cmd_ids[cmd]
                    # skip commands with high ids
                    if num_cmds is not None and cmd_id >= num_cmds:
                        continue

                    x_values.append(day_no)
                    y_values.append(cmd_id)
                    s_values.append(count)
                    c_values.append(tab_colors[user_idx])

                today_cmds = defaultdict(int)

                this_day += SECONDS_IN_A_DAY
                day_no += 1
                while timestamp >= this_day + SECONDS_IN_A_DAY:
                    this_day += SECONDS_IN_A_DAY
                    day_no += 1

            if cmd not in cmd_ids:
                cmd_ids[cmd] = len(cmd_ids)
                y_labels.append(matplotlib_escape(cmd))

            today_cmds[cmd] += 1

        all_x_values.extend(x_values)
        all_y_values.extend(y_values)
        all_s_values.extend(s_values)
        all_c_values.extend(c_values)
        x_values = []
        y_values = []
        s_values = []
        c_values = []
        legend_patches.append(mpatches.Patch(color=tab_colors[user_idx], label="{} ({}) (TODO: sanitize!)".format(name, user_idx)))

    if num_cmds is not None and len(y_labels) > num_cmds:
        y_labels = y_labels[:num_cmds]
    plt.yticks(ticks=range(0, len(y_labels)), labels=y_labels, fontsize=6)
    plt.scatter(all_x_values, all_y_values, s=all_s_values, c=all_c_values, marker='o')
    plt.legend(handles=legend_patches, loc="best")

    if async_draw:
        plt.draw()
    else:
        plt.show()


# Figure 5.6. Command line vocabulary size vs. the number of commands entered for four typical individuals.
def plot_cmdVocabularySize_time():
    SECONDS_IN_A_DAY = 86400
    plt.figure(figsize=(PLOT_WIDTH, PLOT_HEIGHT))
    plt.title("Command vocabulary size growth in time")
    plt.ylabel("Command vocabulary size")
    plt.xlabel("Days")
    legend = []

    # x_count = max(map(lambda x: len(x[1]), DATA_records_by_user.items()))
    # x_values = range(0, x_count)  
    for user in DATA_records_by_user.items():
        new_cmds_after_100 = 0
        new_cmds_after_200 = 0
        new_cmds_after_300 = 0
        cmd_vocabulary = set()
        y_cmd_count = [0]
        name, records = user

        cmd_fail_count = 0

        if not len(records):
            print("ERROR: no records for user {}".format(name))
            continue

        first_day = records[0]["realtimeAfter"]
        this_day = first_day

        for record in records:
            cmd = record["command"]
            timestamp = record["realtimeAfter"]

            if cmd == "":
                cmd_fail_count += 1
                continue

            if timestamp >= this_day + SECONDS_IN_A_DAY:
                this_day += SECONDS_IN_A_DAY
                while timestamp >= this_day + SECONDS_IN_A_DAY:
                    y_cmd_count.append(-10)
                    this_day += SECONDS_IN_A_DAY

                y_cmd_count.append(len(cmd_vocabulary))

                if len(y_cmd_count) > 100:
                    new_cmds_after_100+=1
                if len(y_cmd_count) > 200:
                    new_cmds_after_200+=1
                if len(y_cmd_count) > 300:
                    new_cmds_after_300+=1

                if len(y_cmd_count) == 100:
                    print("% {}: Cmd adoption rate at 100 days (between 0 and 100 days) = {}".format(name, len(cmd_vocabulary) / (len(y_cmd_count))))
                if len(y_cmd_count) == 200:
                    print("% {}: Cmd adoption rate at 200 days days = {}".format(name, len(cmd_vocabulary) / (len(y_cmd_count))))
                    print("% {}: Cmd adoption rate between 100 and 200 days = {}".format(name, new_cmds_after_100 / (len(y_cmd_count) - 100)))
                if len(y_cmd_count) == 300:
                    print("% {}: Cmd adoption rate between 200 and 300 days = {}".format(name, new_cmds_after_200 / (len(y_cmd_count) - 200)))

            if cmd not in cmd_vocabulary:
                cmd_vocabulary.add(cmd)  
        

        print("% {}: New cmd adoption rate after 100 days = {}".format(name, new_cmds_after_100 / (len(y_cmd_count) - 100)))
        print("% {}: New cmd adoption rate after 200 days = {}".format(name, new_cmds_after_200 / (len(y_cmd_count) - 200)))
        print("% {}: New cmd adoption rate after 300 days = {}".format(name, new_cmds_after_300 / (len(y_cmd_count) - 300)))
        print("% {}: cmd_fail_count = {}".format(name, cmd_fail_count))
        x_cmds_entered = range(0, len(y_cmd_count))
        plt.plot(x_cmds_entered, y_cmd_count, 'o', markersize=2)
        legend.append(name + " (TODO: sanitize!)")

    # print(cmd_vocabulary)

    plt.legend(legend, loc="best")
    plt.ylim(bottom=0)

    if async_draw:
        plt.draw()
    else:
        plt.show()


# Figure 5.6. Command line vocabulary size vs. the number of commands entered for four typical individuals.
def plot_cmdLineVocabularySize_cmdLinesEntered():
    plt.figure(figsize=(PLOT_WIDTH, PLOT_HEIGHT))
    plt.title("Command line vocabulary size vs. the number of command lines entered")
    plt.ylabel("Command line vocabulary size")
    plt.xlabel("# of command lines entered")
    legend = []

    for user in DATA_records_by_user.items():
        cmdLine_vocabulary = set()
        y_cmdLine_count = [0]
        name, records = user
        for record in records:
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
        plt.plot(x_cmdLines_entered, y_cmdLine_count, '-')
        legend.append(name + " (TODO: sanitize!)")

    plt.legend(legend, loc="best")

    if async_draw:
        plt.draw()
    else:
        plt.show()

# Figure 3.3. Sequential structure of UNIX command usage, from Figure 4 in Hanson et al. (1984).
#       Ball diameters are proportional to stationary probability. Lines indicate significant dependencies,
#       solid ones being more probable (p < .0001) and dashed ones less probable (.005 < p < .0001).
def graph_cmdSequences(node_count=33, edge_minValue=0.05, view_graph=True):
    START_CMD = "_start_"
    END_CMD = "_end_"
    cmd_count = defaultdict(int)
    cmdSeq_count = defaultdict(lambda: defaultdict(int))
    cmd_id = dict()
    x = 0
    cmd_id[START_CMD] = str(x) 
    x += 1
    cmd_id[END_CMD] = str(x) 
    for pid, session in DATA_records_by_session.items():
        cmd_count[START_CMD] += 1
        prev_cmd = START_CMD
        for record in session:
            cmd = record["command"]
            if cmd == "":
                continue
            cmdSeq_count[prev_cmd][cmd] += 1
            cmd_count[cmd] += 1
            if cmd not in cmd_id:
                x += 1
                cmd_id[cmd] = str(x)
            prev_cmd = cmd
        # end the session
        cmdSeq_count[prev_cmd][END_CMD] += 1
        cmd_count[END_CMD] += 1
        

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
                # if scale_ > 0.5:
                #     graph.edge(cmd_id[cmd], cmd_id[cmd2], constraint='true', splines='curved',
                #             penwidth=penwidth_, style='bold', arrowhead='diamond')
                # elif scale_ > 0.2:
                if scale_ > 0.3:
                    scale_ = str(int(scale_ * 100)/100)
                    graph.edge(cmd_id[cmd], cmd_id[cmd2], constraint='true', splines='curved',
                            penwidth=penwidth_, forcelables='true', label=scale_)
                elif scale_ > 0.2:
                    graph.edge(cmd_id[cmd], cmd_id[cmd2], constraint='true', splines='curved',
                            penwidth=penwidth_, style='dashed')
                # elif scale_ > 0.1:
                else:
                    graph.edge(cmd_id[cmd], cmd_id[cmd2], constraint='false', splines='curved',
                            penwidth=penwidth_, style='dotted', arrowhead='empty')

        # graphviz sometimes fails - see above
        try:
            # graph.view()
            graph.render('/tmp/resh-graph-command_sequence-nodeCount_{}-edgeMinVal_{}.gv'.format(node_count, edge_minValue), view=view_graph)
            break
        except Exception as e:
            trace = traceback.format_exc()
            print("GRAPHVIZ EXCEPTION: <{}>\nGRAPHVIZ TRACE: <{}>".format(str(e), trace))


def plot_strategies_matches(plot_size=50, selected_strategies=[], show_strat_title=True, force_strat_title=None):
    plt.figure(figsize=(PLOT_WIDTH, PLOT_HEIGHT))
    plt.title("Matches at distance <{}>".format(datetime.now().strftime('%H:%M:%S')))
    plt.ylabel('%' + " of matches")
    plt.xlabel("Distance")
    legend = []
    x_values = range(1, plot_size+1)
    saved_matches_total = None
    saved_dataPoint_count = None
    for strategy in data["Strategies"]:
        strategy_title = strategy["Title"]
        # strategy_description = strategy["Description"]

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
            
        # recent is very simple strategy so we will believe 
        #       that there is no bug in it and we can use it to determine total
        if strategy_title == "recent":
            saved_matches_total = matches_total
            saved_dataPoint_count = dataPoint_count

        if len(selected_strategies) and strategy_title not in selected_strategies:
            continue

        acc = 0
        matches_cumulative = []
        for x in matches:
            acc += x
            matches_cumulative.append(acc)
        # matches_cumulative.append(matches_total)
        matches_percent = list(map(lambda x: 100 * x / dataPoint_count, matches_cumulative))

        plt.plot(x_values, matches_percent, 'o-')
        if force_strat_title is not None:
            legend.append(force_strat_title)
        else:
            legend.append(strategy_title)


    assert(saved_matches_total is not None)
    assert(saved_dataPoint_count is not None)
    max_values = [100 * saved_matches_total / saved_dataPoint_count] * len(x_values)
    print("% >>> Avg recurrence rate = {}".format(max_values[0]))
    plt.plot(x_values, max_values, 'r-')
    legend.append("maximum possible")

    x_ticks = list(range(1, plot_size+1, 2))
    x_labels = x_ticks[:]
    plt.xticks(x_ticks, x_labels)
    plt.ylim(bottom=0)
    if show_strat_title:
        plt.legend(legend, loc="best")
    if async_draw:
        plt.draw()
    else:
        plt.show()


def plot_strategies_charsRecalled(plot_size=50, selected_strategies=[]):
    plt.figure(figsize=(PLOT_WIDTH, PLOT_HEIGHT))
    plt.title("Average characters recalled at distance <{}>".format(datetime.now().strftime('%H:%M:%S')))
    plt.ylabel("Average characters recalled")
    plt.xlabel("Distance")
    x_values = range(1, plot_size+1)
    legend = []
    saved_charsRecalled_total = None
    saved_dataPoint_count = None
    for strategy in data["Strategies"]:
        strategy_title = strategy["Title"]
        # strategy_description = strategy["Description"]

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
            
        # recent is very simple strategy so we will believe 
        #       that there is no bug in it and we can use it to determine total
        if strategy_title == "recent":
            saved_charsRecalled_total = charsRecalled_total
            saved_dataPoint_count = dataPoint_count

        if len(selected_strategies) and strategy_title not in selected_strategies:
            continue

        acc = 0
        charsRecalled_cumulative = []
        for x in charsRecalled:
            acc += x
            charsRecalled_cumulative.append(acc)
        charsRecalled_average = list(map(lambda x: x / dataPoint_count, charsRecalled_cumulative))

        plt.plot(x_values, charsRecalled_average, 'o-')
        legend.append(strategy_title)

    assert(saved_charsRecalled_total is not None)
    assert(saved_dataPoint_count is not None)
    max_values = [saved_charsRecalled_total / saved_dataPoint_count] * len(x_values)
    print("% >>> Max avg recalled characters = {}".format(max_values[0]))
    plt.plot(x_values, max_values, 'r-')
    legend.append("maximum possible")

    x_ticks = list(range(1, plot_size+1, 2))
    x_labels = x_ticks[:]
    plt.xticks(x_ticks, x_labels)
    plt.ylim(bottom=0)
    plt.legend(legend, loc="best")
    if async_draw:
        plt.draw()
    else:
        plt.show()


def plot_strategies_charsRecalled_prefix(plot_size=50, selected_strategies=[]):
    plt.figure(figsize=(PLOT_WIDTH, PLOT_HEIGHT))
    plt.title("Average characters recalled at distance (including prefix matches) <{}>".format(datetime.now().strftime('%H:%M:%S'))) 
    plt.ylabel("Average characters recalled (including prefix matches)")
    plt.xlabel("Distance")
    x_values = range(1, plot_size+1)
    legend = []
    saved_charsRecalled_total = None
    saved_dataPoint_count = None
    for strategy in data["Strategies"]:
        strategy_title = strategy["Title"]
        # strategy_description = strategy["Description"]

        dataPoint_count = 0
        matches_total = 0
        charsRecalled = [0] * plot_size
        charsRecalled_total = 0
        
        for multiMatch in strategy["PrefixMatches"]:
            dataPoint_count += 1

            if not multiMatch["Match"]:
                continue
            matches_total += 1

            last_charsRecalled = 0
            for match in multiMatch["Entries"]:

                chars = match["CharsRecalled"]
                charsIncrease = chars - last_charsRecalled
                assert(charsIncrease > 0)
                charsRecalled_total += charsIncrease 

                dist = match["Distance"]  
                if dist > plot_size:
                    continue

                charsRecalled[dist-1] += charsIncrease
                last_charsRecalled = chars
            
        # recent is very simple strategy so we will believe 
        #       that there is no bug in it and we can use it to determine total
        if strategy_title == "recent":
            saved_charsRecalled_total = charsRecalled_total
            saved_dataPoint_count = dataPoint_count

        if len(selected_strategies) and strategy_title not in selected_strategies:
            continue

        acc = 0
        charsRecalled_cumulative = []
        for x in charsRecalled:
            acc += x
            charsRecalled_cumulative.append(acc)
        charsRecalled_average = list(map(lambda x: x / dataPoint_count, charsRecalled_cumulative))

        plt.plot(x_values, charsRecalled_average, 'o-')
        legend.append(strategy_title)

    assert(saved_charsRecalled_total is not None)
    assert(saved_dataPoint_count is not None)
    max_values = [saved_charsRecalled_total / saved_dataPoint_count] * len(x_values)
    print("% >>> Max avg recalled characters (including prefix matches) = {}".format(max_values[0]))
    plt.plot(x_values, max_values, 'r-')
    legend.append("maximum possible")

    x_ticks = list(range(1, plot_size+1, 2))
    x_labels = x_ticks[:]
    plt.xticks(x_ticks, x_labels)
    plt.ylim(bottom=0)
    plt.legend(legend, loc="best")
    if async_draw:
        plt.draw()
    else:
        plt.show()


def plot_strategies_matches_noncummulative(plot_size=50, selected_strategies=["recent (bash-like)"], show_strat_title=False, force_strat_title=None):
    plt.figure(figsize=(PLOT_WIDTH, PLOT_HEIGHT))
    plt.title("Matches at distance (noncumulative) <{}>".format(datetime.now().strftime('%H:%M:%S')))
    plt.ylabel('%' + " of matches")
    plt.xlabel("Distance")
    legend = []
    x_values = range(1, plot_size+1)
    saved_matches_total = None
    saved_dataPoint_count = None
    for strategy in data["Strategies"]:
        strategy_title = strategy["Title"]
        # strategy_description = strategy["Description"]

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
            
        # recent is very simple strategy so we will believe 
        #       that there is no bug in it and we can use it to determine total
        if strategy_title == "recent":
            saved_matches_total = matches_total
            saved_dataPoint_count = dataPoint_count

        if len(selected_strategies) and strategy_title not in selected_strategies:
            continue

        # acc = 0
        # matches_cumulative = []
        # for x in matches:
        #     acc += x
        #     matches_cumulative.append(acc)
        # # matches_cumulative.append(matches_total)
        matches_percent = list(map(lambda x: 100 * x / dataPoint_count, matches))

        plt.plot(x_values, matches_percent, 'o-')
        if force_strat_title is not None:
            legend.append(force_strat_title)
        else:
            legend.append(strategy_title)

    assert(saved_matches_total is not None)
    assert(saved_dataPoint_count is not None)
    # max_values = [100 * saved_matches_total / saved_dataPoint_count] * len(x_values)
    # print("% >>> Avg recurrence rate = {}".format(max_values[0]))
    # plt.plot(x_values, max_values, 'r-')
    # legend.append("maximum possible")

    x_ticks = list(range(1, plot_size+1, 2))
    x_labels = x_ticks[:]
    plt.xticks(x_ticks, x_labels)
    # plt.ylim(bottom=0)
    if show_strat_title:
        plt.legend(legend, loc="best")
    if async_draw:
        plt.draw()
    else:
        plt.show()


def plot_strategies_charsRecalled_noncummulative(plot_size=50, selected_strategies=["recent (bash-like)"], show_strat_title=False):
    plt.figure(figsize=(PLOT_WIDTH, PLOT_HEIGHT))
    plt.title("Average characters recalled at distance (noncumulative) <{}>".format(datetime.now().strftime('%H:%M:%S')))
    plt.ylabel("Average characters recalled")
    plt.xlabel("Distance")
    x_values = range(1, plot_size+1)
    legend = []
    saved_charsRecalled_total = None
    saved_dataPoint_count = None
    for strategy in data["Strategies"]:
        strategy_title = strategy["Title"]
        # strategy_description = strategy["Description"]

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
            
        # recent is very simple strategy so we will believe 
        #       that there is no bug in it and we can use it to determine total
        if strategy_title == "recent":
            saved_charsRecalled_total = charsRecalled_total
            saved_dataPoint_count = dataPoint_count

        if len(selected_strategies) and strategy_title not in selected_strategies:
            continue

        # acc = 0
        # charsRecalled_cumulative = []
        # for x in charsRecalled:
        #     acc += x
        #     charsRecalled_cumulative.append(acc)
        # charsRecalled_average = list(map(lambda x: x / dataPoint_count, charsRecalled_cumulative))
        charsRecalled_average = list(map(lambda x: x / dataPoint_count, charsRecalled))

        plt.plot(x_values, charsRecalled_average, 'o-')
        legend.append(strategy_title)

    assert(saved_charsRecalled_total is not None)
    assert(saved_dataPoint_count is not None)
    # max_values = [saved_charsRecalled_total / saved_dataPoint_count] * len(x_values)
    # print("% >>> Max avg recalled characters = {}".format(max_values[0]))
    # plt.plot(x_values, max_values, 'r-')
    # legend.append("maximum possible")

    x_ticks = list(range(1, plot_size+1, 2))
    x_labels = x_ticks[:]
    plt.xticks(x_ticks, x_labels)
    # plt.ylim(bottom=0)
    if show_strat_title:
        plt.legend(legend, loc="best")
    if async_draw:
        plt.draw()
    else:
        plt.show()


def plot_strategies_charsRecalled_prefix_noncummulative(plot_size=50, selected_strategies=["recent (bash-like)"], show_strat_title=False):
    plt.figure(figsize=(PLOT_WIDTH, PLOT_HEIGHT))
    plt.title("Average characters recalled at distance (including prefix matches) (noncummulative) <{}>".format(datetime.now().strftime('%H:%M:%S'))) 
    plt.ylabel("Average characters recalled (including prefix matches)")
    plt.xlabel("Distance")
    x_values = range(1, plot_size+1)
    legend = []
    saved_charsRecalled_total = None
    saved_dataPoint_count = None
    for strategy in data["Strategies"]:
        strategy_title = strategy["Title"]
        # strategy_description = strategy["Description"]

        dataPoint_count = 0
        matches_total = 0
        charsRecalled = [0] * plot_size
        charsRecalled_total = 0
        
        for multiMatch in strategy["PrefixMatches"]:
            dataPoint_count += 1

            if not multiMatch["Match"]:
                continue
            matches_total += 1

            last_charsRecalled = 0
            for match in multiMatch["Entries"]:

                chars = match["CharsRecalled"]
                charsIncrease = chars - last_charsRecalled
                assert(charsIncrease > 0)
                charsRecalled_total += charsIncrease 

                dist = match["Distance"]  
                if dist > plot_size:
                    continue

                charsRecalled[dist-1] += charsIncrease
                last_charsRecalled = chars
            
        # recent is very simple strategy so we will believe 
        #       that there is no bug in it and we can use it to determine total
        if strategy_title == "recent":
            saved_charsRecalled_total = charsRecalled_total
            saved_dataPoint_count = dataPoint_count

        if len(selected_strategies) and strategy_title not in selected_strategies:
            continue

        # acc = 0
        # charsRecalled_cumulative = []
        # for x in charsRecalled:
        #     acc += x
        #     charsRecalled_cumulative.append(acc)
        # charsRecalled_average = list(map(lambda x: x / dataPoint_count, charsRecalled_cumulative))
        charsRecalled_average = list(map(lambda x: x / dataPoint_count, charsRecalled))

        plt.plot(x_values, charsRecalled_average, 'o-')
        legend.append(strategy_title)

    assert(saved_charsRecalled_total is not None)
    assert(saved_dataPoint_count is not None)
    # max_values = [saved_charsRecalled_total / saved_dataPoint_count] * len(x_values)
    # print("% >>> Max avg recalled characters (including prefix matches) = {}".format(max_values[0]))
    # plt.plot(x_values, max_values, 'r-')
    # legend.append("maximum possible")

    x_ticks = list(range(1, plot_size+1, 2))
    x_labels = x_ticks[:]
    plt.xticks(x_ticks, x_labels)
    # plt.ylim(bottom=0)
    if show_strat_title:
        plt.legend(legend, loc="best")
    if async_draw:
        plt.draw()
    else:
        plt.show()


def print_top_cmds(num_cmds=20):
    cmd_count = defaultdict(int)
    cmd_total = 0
    for pid, session in DATA_records_by_session.items():
        for record in session:
            cmd = record["command"]
            if cmd == "":
                continue
            cmd_count[cmd] += 1
            cmd_total += 1

    # get `node_count` of largest nodes
    sorted_cmd_count = list(sorted(cmd_count.items(), key=lambda x: x[1], reverse=True))
    print("\n\n% All subjects: Top commands")
    for cmd, count in sorted_cmd_count[:num_cmds]:
        print("{} {}".format(cmd, count))
    # print(sorted_cmd_count)
    # cmds_to_graph = list(map(lambda x: x[0], sorted_cmd_count))[:cmd_count]


def print_top_cmds_by_user(num_cmds=20):
    for user in DATA_records_by_user.items():
        name, records = user
        cmd_count = defaultdict(int)
        cmd_total = 0
        for record in records:
            cmd = record["command"]
            if cmd == "":
                continue
            cmd_count[cmd] += 1
            cmd_total += 1

        # get `node_count` of largest nodes
        sorted_cmd_count = list(sorted(cmd_count.items(), key=lambda x: x[1], reverse=True))
        print("\n\n% {}: Top commands".format(name))
        for cmd, count in sorted_cmd_count[:num_cmds]:
            print("{} {}".format(cmd, count))
        # print(sorted_cmd_count)
        # cmds_to_graph = list(map(lambda x: x[0], sorted_cmd_count))[:cmd_count]


def print_avg_cmdline_length():
    cmd_len_total = 0
    cmd_total = 0
    for pid, session in DATA_records_by_session.items():
        for record in session:
            cmd = record["cmdLine"]
            if cmd == "":
                continue
            cmd_len_total += len(cmd) 
            cmd_total += 1

    print("% ALL avg cmdline = {}".format(cmd_len_total / cmd_total))
    # print(sorted_cmd_count)
    # cmds_to_graph = list(map(lambda x: x[0], sorted_cmd_count))[:cmd_count]


# plot_cmdLineFrq_rank()
# plot_cmdFrq_rank()
print_top_cmds(30)
print_top_cmds_by_user(30)
# print_avg_cmdline_length()
#         
# plot_cmdLineVocabularySize_cmdLinesEntered()
plot_cmdVocabularySize_cmdLinesEntered()
plot_cmdVocabularySize_time()
# plot_cmdVocabularySize_daily()
plot_cmdUsage_in_time(num_cmds=100)
plot_cmdUsage_in_time(sort_cmds=True, num_cmds=100)
# 
recent_strats=("recent", "recent (bash-like)")
recurrence_strat=("recent (bash-like)",)
# plot_strategies_matches(20)
# plot_strategies_charsRecalled(20)
# plot_strategies_charsRecalled_prefix(20)
# plot_strategies_charsRecalled_noncummulative(20, selected_strategies=recent_strats)
# plot_strategies_matches_noncummulative(20)
# plot_strategies_charsRecalled_noncummulative(20)
# plot_strategies_charsRecalled_prefix_noncummulative(20)
# plot_strategies_matches(20, selected_strategies=recurrence_strat, show_strat_title=True, force_strat_title="recurrence rate")
# plot_strategies_matches_noncummulative(20, selected_strategies=recurrence_strat, show_strat_title=True, force_strat_title="recurrence rate")

# graph_cmdSequences(node_count=33, edge_minValue=0.048)

# graph_cmdSequences(node_count=28, edge_minValue=0.06)

# new improved
# for n in range(40, 43):
#     for e in range(94, 106, 2):
#         e *= 0.001
#         graph_cmdSequences(node_count=n, edge_minValue=e, view_graph=False)

#for n in range(29, 35):
#    for e in range(44, 56, 2):
#        e *= 0.001
#        graph_cmdSequences(node_count=n, edge_minValue=e, view_graph=False)

# be careful and check if labels fit the display

if async_draw:
    plt.show()
