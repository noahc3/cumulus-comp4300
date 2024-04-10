#include <iostream>
#include <string>

const uint64_t SEED = 0x12345678;

int main(int argc, char** argv) {
    std::string line;

    std::cout << "Type anything, it will be echoed back." << std::endl;

    while (std::getline(std::cin, line)) {
        std::cout << "You typed " << line << std::endl;
    }

    return 0;
}