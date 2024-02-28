#include <iostream>
#include <iomanip> 
#include <string>
#include <fstream>
#include <omp.h>
#include <cstdint>
#include <functional>
#include <sstream>
#include <thread>
#include <vector>
#include "nlohmann/json.hpp"

using json = nlohmann::json;

// Define a City class
class City {
public:
    std::string name;
    int population;
    double area;

    /// <summary>
    /// Parse JSON data and populate the object
    /// </summary>
    /// <param name="dataString">Data to parse to City object</param>
    void fromJson(json dataString) {
        std::string temp = dataString["name"];
        name = temp;
        population = dataString["population"];
        area = dataString["area"];
    }

    /// <summary>
    /// Generate a hashcode for the City object
    /// </summary>
    /// <returns>An int of the city object hash code</returns>
    int hashCode() const {
        std::string data = name + std::to_string(population) + std::to_string(static_cast<int>(area));
        std::hash<std::string> hash_fn;
        return static_cast<int>(hash_fn(data));
    }
};

/// <summary>
/// Define a Monitor class to manage City objects
/// </summary>
class List {
private:
    std::vector<City> list;
    omp_lock_t lock;

public:
    int iSum = 0;
    double dSum = 0;

    List() {
        omp_init_lock(&lock);
    }

    ~List() {
        omp_destroy_lock(&lock);
    }

    /// <summary>
    /// Add a City object to the list in a thread-safe manner
    /// </summary>
    /// <param name="city">City to add</param>
    void Add(City& city) {
        #pragma omp critical
        {
            auto it = std::lower_bound(list.begin(), list.end(), city, [](const City& a, const City& b) {
                return a.population > b.population;
                });
            list.insert(it, city);
        }
    }

    /// <summary>
    /// Remove and return a City object from the list in a thread-safe manner
    /// </summary>
    /// <returns>Returns the removed City object</returns>
    City Pop() {
        City city;

        #pragma omp critical
        {
            if (!list.empty()) {
                city = list.back();
                list.pop_back();
            }
        }
        return city;
    }

    /// <summary>
    /// Gets the count of City objects in the list
    /// </summary>
    /// <returns>Returns the ammount of City objects in the list</returns>
    int Count() const {
        return list.size();
    };
};


/// <summary>
/// Utility class for file I/O
/// </summary>
class IO {
public:
    /// <summary>
    /// Read City data from a JSON file
    /// </summary>
    /// <param name="fileName">file path</param>
    /// <returns>Returns a list of all cities read from the JSON file</returns>
    static std::vector<City> ReadFile(std::string& fileName) {
        std::vector<City> cities;
        std::ifstream stream;
        stream.open(fileName);

        json allCitiesJson = json::parse(stream);
        auto allCities = allCitiesJson["cities"];
        for (const json& city : allCities) {
            City tempCity;
            tempCity.fromJson(city);
            cities.push_back(tempCity);
        }
        stream.close();
        return cities;
    }

    /// <summary>
    /// Print City data and result to a file
    /// </summary>
    /// <param name="fileName">File to write to</param>
    /// <param name="list">List of all cities</param>
    /// <param name="iSum">Sum of all integers</param>
    /// <param name="dSum">Sum of all doubles</param>
    static void printResult(std::string& fileName, List& list, int iSum, double dSum) {
        const char* filePath = fileName.c_str();

        std::ofstream file;
        file.open(filePath, std::ios_base::app);
        file << std::string(30, ' ') <<  "Results" << std::endl;
        file << "+"<< std::string(4, '-') << "+" << std::string(26, '-') << "+" << std::string(20, '-') << "+" << std::string(21, '-') << "+" << std::endl;
        file << "|" << std::setw(5) << "# |" << std::setw(27) << "Name |" << std::setw(21) << "Population |" << std::setw(22) << "Area |" << std::endl;
        file << "+" << std::string(4, '-') << "+" << std::string(26, '-') << "+" << std::string(20, '-') << "+" << std::string(21, '-') << "+" << std::endl;
        int i = 0;
        while (list.Count() > 0) {
            City temp = list.Pop();
            file << "|" << std::setw(3) << i << " |" << std::setw(26) << temp.name << "|" << std::setw(20) << temp.population << "|" << std::setw(20) << temp.area << " |" << std::endl;
            i++;
        }
        file << "+" << std::string(4, '-') << "+" << std::string(26, '-') << "+" << std::string(20, '-') << "+" << std::string(21, '-') << "+" << std::endl;
        file << "Total sum of int fields:" << iSum << std::endl;
        file << "Total sum of double fields:" << dSum << std::endl;
    }

    /// <summary>
    /// Print City data to a file
    /// </summary>
    /// <param name="fileName">File to write to</param>
    /// <param name="cities">List of all files</param>
    static void printResult(std::string& fileName, std::vector<City>& cities) {
        const char* filePath = fileName.c_str();

        remove(filePath);

        std::ofstream file;
        file.open(filePath);
        file << std::string(30, ' ') << "Data" << std::endl;
        file << "+" << std::string(4, '-') << "+" << std::string(20, '-') << "+" << std::string(20, '-') << "+" << std::string(21, '-') << "+" << std::endl;
        file << "|" << std::setw(5) << "# |" << std::setw(21) << "Name |" << std::setw(21) << "Population |" << std::setw(22) << "Area |" << std::endl;
        file << "+" << std::string(4, '-') << "+" << std::string(20, '-') << "+" << std::string(20, '-') << "+" << std::string(21, '-') << "+" << std::endl;
        for (int i = 0; i < cities.size(); i++)
        {
            City temp = cities[i];
            file << "|" << std::setw(3) << i << " |" << std::setw(20) << temp.name << "|" << std::setw(20) << temp.population << "|" << std::setw(20) << temp.area << " |" << std::endl;
        }
        file << "+" << std::string(4, '-') << "+" << std::string(20, '-') << "+" << std::string(20, '-') << "+" << std::string(21, '-') << "+" << std::endl;
        file << std::endl;
    }
};

/// <summary>
/// Utility class for mathematical operations
/// </summary>
class Utils {
public:
    static int fib(int x) {
        if ((x == 1) || (x == 0)) {
            return(x);
        }
        else {
            return(fib(x - 1) + fib(x - 2));
        }
    }
};


/// <summary>
/// Function to process City objects and populate the list
/// </summary>
/// <param name="cities">List of all City objects</param>
/// <param name="list">List to populate the data to</param>
/// <param name="iSum">Sum of all integers</param>
/// <param name="dSum">Sum of all doubles</param>
void execute(std::vector<City>& cities, List& list, int& iSum, double& dSum) {
    for (City& city : cities) {
        int hashCode = city.hashCode();
        int fib = Utils::fib(35);

        if (city.population >= 20000) {
            list.Add(city);
            iSum += city.population;
            dSum += city.area;
        }
    }
}

int main() {
    std::string inputFile = "";
    std::string outputFile = "IFF-1-1_ZekonisE_L1_res.txt";

    int inputNum;
    std::cout << "Choose which file you want to use: ";
    std::cin >> inputNum;

    switch (inputNum) {
        case 1:
            inputFile = "IFF-1-1_ZekonisE_L1_dat_1.json";
            break;
        case 2:
            inputFile = "IFF-1-1_ZekonisE_L1_dat_2.json";
            break;
        case 3:
            inputFile = "IFF-1-1_ZekonisE_L1_dat_3.json";
            break;
            
        default:
            std::cout << "Unknown number provided: " << inputNum << std::endl;
            return 1;
    }

    std::cout << "Input file: " << inputFile << std::endl;


    auto cities = IO::ReadFile(inputFile);
    int numCities = cities.size();
    const int numThreads = std::max(2, static_cast<int>(std::ceil(static_cast<double>(numCities) / 2)));
    std::vector<std::vector<City>> threadData(numThreads);
    
    int citiesInThread = numCities / numThreads;
    int remainingCities = numCities % numThreads;

    int start = 0;

    for (int i = 0; i < numThreads; i++) {
        int end = start + citiesInThread;

        if (remainingCities > 0) {
            end++;
            remainingCities--;
        }

        std::vector<City> batch;
        for (int j = start; j < end; j++) {
            batch.push_back(cities[j]);
        }

        threadData[i] = std::move(batch);

        start = end;
    }

    List monitor;
    int iSum = 0;
    double dSum = 0.00;

    #pragma omp parallel num_threads(numThreads) reduction(+:iSum) reduction(+:dSum)
    {
        int threadID = omp_get_thread_num();
        execute(threadData[threadID], monitor, iSum, dSum);
    }

    // Print results to the output file
    IO::printResult(outputFile, cities);
    IO::printResult(outputFile, monitor, iSum, dSum);

    return 0;
}