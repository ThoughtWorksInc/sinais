names = []

with open('names.txt') as fp:
    for line in fp:
        pt, en = line.split('\t')
        names.append((pt, en.strip()))

def pt_len(pair):
    return len(pair[0])

names = sorted(names, key=pt_len, reverse=True)

srcs = 'runefinder.go runefinder_test.go'.split()

for src in srcs:
    src = '../' + src
    with open(src) as fp:
        txt = fp.read()
        for pt, en in names:
            txt = txt.replace(pt, en)
    with open(src, 'wt') as fp:
        fp.write(txt)
