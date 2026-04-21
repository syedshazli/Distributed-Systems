package mapreduce

import (
	"encoding/json"
	"os"
	"sort"
)

// doReduce does the job of a reduce worker: it reads the intermediate
// key/value pairs (produced by the map phase) for this task, sorts the
// intermediate key/value pairs by key, calls the user-defined reduce function
// (reduceF) for each key, and writes the output to disk.
func doReduce(
	jobName string, // the name of the whole MapReduce job
	reduceTaskNumber int, // which reduce task this is
	nMap int, // the number of map tasks that were run ("M" in the paper)
	reduceF func(key string, values []string) string,
) {
	// TODO:
	// You will need to write this function.
	//
	// Step 1: Read all intermediate files for this reduce task.
	//   There are nMap intermediate files, one per map task.
	//   For map task m, the file is: reduceName(jobName, m, reduceTaskNumber)
	//   Open each file, create a json.NewDecoder, and call .Decode() in a loop
	//   until it returns an error (io.EOF means you've read all entries):
	//     dec := json.NewDecoder(f)
	//     for {
	//       var kv KeyValue // Assuming KeyValue is {Key string, Value string}
	//       if err := dec.Decode(&kv); err != nil { break } // io.EOF comes here

	//     }
	//

	// Use a map to aggregate all values for each key across all intermediate files.
	intermediateData := make(map[string][]string)

	for fileNum := 0; fileNum < nMap; fileNum++ {
		fileName := reduceName(jobName, fileNum, reduceTaskNumber)

		file, err0 := os.Open(fileName)
		checkError(err0)

		decoder := json.NewDecoder(file)

		for {
			var kv KeyValue
			if err := decoder.Decode(&kv); err != nil {
				break // End of file or error
			}
			intermediateData[kv.Key] = append(intermediateData[kv.Key], kv.Value)
		}
		file.Close()
	}

	// Step 2: Call reduceF for each unique key.
	//   Collect all keys, sort them (sort.Strings), then for each key call
	//   reduceF(key, kvMap[key]) to get the output value.
	var keys []string
	for key := range intermediateData {
		keys = append(keys, key)
	}
	sort.Strings(keys) // Sort the unique keys

	// Step 3: Write output atomically using a temp file + os.Rename.
	//   The final output filename is: mergeName(jobName, reduceTaskNumber)
	//   Use os.CreateTemp to create a temporary output file.
	//   Write JSON-encoded KeyValue pairs to this temporary file.
	//   Close the temporary file, then call os.Rename(tmp.Name(), mergeName(...)).
	finalOutputFileName := mergeName(jobName, reduceTaskNumber)
	tmpFile, err := os.CreateTemp(".", "mr-reduce-output-") // Use a descriptive prefix
	checkError(err)

	enc := json.NewEncoder(tmpFile)
	for _, key := range keys {
		sort.Strings(intermediateData[key]) // Sort values for the current key
		reducedValue := reduceF(key, intermediateData[key])
		err := enc.Encode(KeyValue{key, reducedValue})
		checkError(err)
	}

	err = tmpFile.Close()
	checkError(err)

	os.Remove(finalOutputFileName)
	err = os.Rename(tmpFile.Name(), finalOutputFileName)
	checkError(err)
}
