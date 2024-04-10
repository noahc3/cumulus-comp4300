# Variables
CXX = g++
CXXFLAGS = -Wall -std=c++11
TARGET = mock.out
SRC = mock.cpp
TARGET2 = echo.out
SRC2 = echo.cpp

# Default rule
all: $(TARGET) $(TARGET2)

# Rule to build the target
$(TARGET): $(SRC)
	$(CXX) $(CXXFLAGS) -o $(TARGET) $(SRC)

# Rule to build the second target
$(TARGET2): $(SRC2)
	$(CXX) $(CXXFLAGS) -o $(TARGET2) $(SRC2)

# Rule to clean the workspace
clean:
	rm -f $(TARGET) $(TARGET2)
