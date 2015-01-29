totalf = 0
with open('sbout.txt') as f:
    for d in f:
			itemName = d.split(" with frequency ")[0]
			freq = d.split(" with frequency ")[1].split(" at least ")[0]
			print("ItemName: " + itemName + " freq: " + freq)	
			totalf += int(freq)

print(totalf)