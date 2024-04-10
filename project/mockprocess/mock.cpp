#include <iostream>
#include <fstream>
#include <string>
#include <random>
#include <thread>
#include <chrono>

const uint64_t SEED = 0x12345678;

int main(int argc, char** argv) {
    if (argc < 4 || argc > 4) {
        std::cout << "Produce mock output on stdout for developing applications which read standard output." << std::endl;
        std::cout << "When the end of the mock data file is reached, the program repeats from the beginning." << std::endl;
        std::cout << "Press CTRL+C to exit." << std::endl;
        std::cout << "Usage: " << argv[0] << " <filename> <min lines per second> <max lines per second>" << std::endl;
    }

    std::string filename = argv[1];
    int minlps = std::stoi(argv[2]);
    int maxlps = std::stoi(argv[3]);

    float mindelay = 1.0 / minlps;
    float maxdelay = 1.0 / maxlps;

    std::mt19937 gen(SEED);
    std::uniform_real_distribution<> dis(mindelay, maxdelay);

    std::ifstream file(filename);
    if (file.is_open()) {
        std::string line;

        while (file.is_open()) {
            while(std::getline(file, line)) {
                float delay = dis(gen);
                std::cout << line << std::endl;
                std::this_thread::sleep_for(std::chrono::milliseconds((int)(delay * 1000)));
            }

            file.seekg(0, std::ios::beg);
        }
    } else {
        std::cout << "Failed to open file '" << filename << "' for reading." << std::endl;
    }




    return 0;
}