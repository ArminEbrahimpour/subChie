uniq_lines = set(open('output.txt').readlines())

uniq_output = open('output.txt', 'w').writelines(uniq_lines)


import re

text = open('output.txt').readlines()

output = ""

for i in text:
    # it is a list of urls which starts with accounts.google...
    lines_starting_with_some_words = re.findall(r"^https://accounts.google.com/*", i, flags=re.IGNORECASE | re.MULTILINE)

    # if it's a valid url write it in output
    if lines_starting_with_some_words != ['https://accounts.google.com/']:
        print(i)
        output += i+'\n'

# write final output to the output.txt file
with open('output.txt', 'w') as file:
    file.write(output)