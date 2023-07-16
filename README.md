# MD5-Hasher
Hashing the response in MD5 format

This tool makes http requests and prints the address of the
request along with the MD5 hash of the response.

example : 
$ ./myapp -parallel 5 http://example.com http://google.com http://github.com

Address: http://github.com, Hash: c07b3f8381c1971d9793e70c7be64f5b

Address: http://example.com, Hash: 84238dfc8092e5d9c0dac8ef93371a07

Address: http://google.com, Hash: b90313d4f39eba83754d0f4c6a8dd852 

Useage : The tool takes the target address and parallel flag. 
        Note: The parallel flag provides the maximum limit of parallel requests,, with is 10 as defalut
