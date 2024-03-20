#include <iostream>
#include <fstream>
#include <string>

int main() {
    std::ifstream inputFile("/source.txt"); // Open the file
    if (!inputFile) { // Check if the file opened successfully
        std::cerr << "Unable to open the file!" << std::endl;
        return 1;
    }

    std::string line;
    long long lines = 0;
    while (std::getline(inputFile, line)) { // Read the file line by line
        int n = 0;
        while (n++ < 30000) {}
        if (lines % 10000 == 0) {
            std::cout << lines / 10000 << std::endl;
        }
        lines++;
    }

    inputFile.close(); // Close the file
	
    std::cout << "lines: " << lines;
    return 0;
}
