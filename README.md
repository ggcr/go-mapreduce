# MapReduce

## Background
MapReduce comes from a paper by two actual Google Senior Fellows, Jeff Dean and Sanjay Ghemawat. They both are the only awarded with this title which is the highest position at Google and one of the main reasons is because they saved the entire company and changed the course of our lives forever.

Back in 2004, Google’s co-founders, Larry Page and Sergey Brin, were negotiating a deal to power a search engine for Yahoo, and they promised to deliver an index ten times bigger than the one they had at the time. What Yahoo did not know is that Google infrastructure wasn’t composed of expensive racks of computers, it was made up of cheap spare parts and it was not working. 

At the time Google was delivering an indexed results outdated by 5 months because their infrastructure was acting as a bottleneck at the moment of computate all the indexes. Google was bigger than its expectations.

[Read more](https://www.newyorker.com/magazine/2018/12/10/the-friendship-that-made-google-huge)
## The need for a Distributed System was born
If we stick to the definition, a Distributed System is a set of cooperating computers that are communicating with each other over the network to get some coherent task done.

But this comes at a high cost and as a solid proof of this we know that it is one of the most active research areas in recent years, now eclipsed by Machine Learning but still present than ever both in the industry and academia. 

The main goal of a Distributed System, aside of course of Performance and Scalability is to deliver an Interface that facilitates the writing of software that will end up at a Distributed System. That’s the main objective.

Google engineers were reaching this goal, to deliver an interface for every engineer (whether they knew a lot about the infrastructure or not) could extend and add features. It was critical for the surveillance of the company.

That’s what MapReduce is all about: an interface.
## The problem
It assumes we have an input split up into M files. 

A map function is a Higher Order Function, that is, it accepts functions as parameters and this function passed as argument will be applied to each input independent chunk.

It takes as input an Input File, applies the desired function and produces a list of key-value pairs as output.
-  **Map function definition**

   As a starter point we can assign our map function to count the occurrences of words in M input files:

- **Reduce function definition**
   
  The reduce function will collect together all the map outputs with the same keys and produce a reduce result (that is shorter than at least the initial  map occurrences).
  
![Captura de Pantalla 2022-10-29 a las 21 00 33](https://user-images.githubusercontent.com/57730982/198848456-c149a612-274d-4c10-ba73-945f01cd69f2.png)

## Terminology

We define the entire Job of MapReduce and then tasks such as a single Map function acting for a single input.

## Implementation
Remember and don’t forget that MapReduce is an interface, the high level programmer won’t need to know any of the previous explanations.

**Map(key, value)** – This is the map function call, key is the filename, and value is the filename content.
```
Split v into words
For each word:
	Emit(word, “1”)
```

**Reduce(key, value)** – This is the reduce function call, key is the word, and value is all the array of “1” we have, we only care on the count.
```
  Emit(word, length(value))
```
  
## Infrastructure

We have a Master process and multiple Workers, the master will indicate each worker to start with the mapping phase (read input, count occurrences and store in intermediate servers local disks). 

So we can see that the “Emit” function is of the workers nature and this function will be to write the Map Output to intermediate servers local disks.
When the master has finished Mapping all the data then it will start communicating to the workers to do the Reduce phase and generate the output files.
![Captura de Pantalla 2022-10-29 a las 21 01 19](https://user-images.githubusercontent.com/57730982/198848484-4b7fb2b5-436b-4dc1-a13c-2855c8a07d55.png)

## The bottleneck of remote procedures
Enough talking, I already convinced you of all the pros and cool things of MapReduce but what are the disadvantages?

The main problem with MapReduce is in the remote Input, if a worker has to read an input file which is not located on the same machine then this will be as fast as the network allows us.

What Google did was to make the Master know the location of the Input splits and generate a worker in each location needed so that the Map Phase was local read and local write. 

The other problem Google encountered was that each worker would store a row of occurrences, and what reduce really needs as the values is the vectorized columns, not the rows:
 ![Captura de Pantalla 2022-10-29 a las 21 01 50](https://user-images.githubusercontent.com/57730982/198848508-64caeebf-8e3c-4a0d-b543-b359944b3ebc.png)
The operation of transforming the Rows into Columns is called Shuffle in the paper.

## How To Run
```bash
$ go run main.go

============== MAP PHASE ==============
file input/pg-being_ernest.txt read ✅
file input/pg-metamorphosis.txt read ✅
file input/pg-huckleberry_finn.txt read ✅
file input/pg-dorian_gray.txt read ✅
file input/pg-tom_sawyer.txt read ✅
file input/pg-frankenstein.txt read ✅
file input/pg-sherlock_holmes.txt read ✅
file tmp/pg-being_ernest.txt wrote 📕
file tmp/pg-metamorphosis.txt wrote 📕
input/pg-being_ernest.txt DONE
input/pg-metamorphosis.txt DONE
file input/pg-grimm.txt read ✅
file tmp/pg-dorian_gray.txt wrote 📕
input/pg-dorian_gray.txt DONE
file tmp/pg-frankenstein.txt wrote 📕
input/pg-frankenstein.txt DONE
file tmp/pg-tom_sawyer.txt wrote 📕
input/pg-tom_sawyer.txt DONE
file tmp/pg-grimm.txt wrote 📕
input/pg-grimm.txt DONE
file tmp/pg-huckleberry_finn.txt wrote 📕
input/pg-huckleberry_finn.txt DONE
file tmp/pg-sherlock_holmes.txt wrote 📕
input/pg-sherlock_holmes.txt DONE

======================== REDUCE PHASE =====================
length for readed dict (Key="tmp/pg-being_ernest.txt") => 3012
length of final map: 3012
length for readed dict (Key="tmp/pg-dorian_gray.txt") => 7119
length of final map: 7977
length for readed dict (Key="tmp/pg-frankenstein.txt") => 7262
length of final map: 11457
length for readed dict (Key="tmp/pg-grimm.txt") => 5134
length of final map: 12880
length for readed dict (Key="tmp/pg-huckleberry_finn.txt") => 6483
length for readed dict (Key="tmp/pg-metamorphosis.txt") => 2982
length of final map: 15516
length of final map: 15817
length for readed dict (Key="tmp/pg-sherlock_holmes.txt") => 8070
length of final map: 17905
length for readed dict (Key="tmp/pg-tom_sawyer.txt") => 7625

length of final map output.json: 19436
107.629025ms
```

## References

J. Dean and S. Ghemawat, “MapReduce: Simplified data processing on large clusters” Jan. 2008.
