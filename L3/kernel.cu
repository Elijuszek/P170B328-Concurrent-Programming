#define __CUDACC__
#include "cuda_runtime.h"
#include "device_launch_parameters.h"
#include <cuda.h>
#include <iostream>
#include <iomanip>
#include <stdio.h>
#include <fstream>

#include "nlohmann/json.hpp"

using namespace std;
using json = nlohmann::json;

class City {
public:
    char name[256];
    int population;
    double area;
    char res[256];
};

//const string dataFile = "IFF-1-1_ZekonisE_L1_dat_1.json"; // visi tinka
const string dataFile = "IFF-1-1_ZekonisE_L1_dat_2.json"; // kaikurie tinka
//const string dataFile = "IFF-1-1_ZekonisE_L1_dat_3.json"; // nei vienas netinka
const string resultFile = "IFF-1-1_ZekonisE_L1_res.txt";

void readCitiesFile(vector<City>* cities) {
    ifstream stream(dataFile);
    json allCitiesJson = json::parse(stream);

    auto allItems = allCitiesJson["cities"];
    for (auto& new_items : allItems) {
        City tempItem;
        string n = new_items["name"];
        //Returns a pointer to an array that contains a null-terminated sequence of characters (i.e., a C-string) 
        // representing the current value of the string object.
        strcpy(tempItem.name, n.c_str());
        tempItem.population = new_items["population"];
        tempItem.area = new_items["area"];
        cities->push_back(tempItem);
    }
    stream.close();
}

void writeListToFile(vector<City>& cities, string fileName) {
    ofstream file;
    file.open(fileName, ios::out);
    file << setw(33) << "Pradiniai duomenys" << endl
        << "--------------------------------------------------------------" << endl
        << setw(5) << "Nr. |" << setw(30) << "Name |" << setw(15) << "Population |" << setw(17) << "Area |" << endl
        << "--------------------------------------------------------------" << endl;

    for (int i = 0; i < cities.size(); i++)
    {
        file << setw(5) << to_string(i+1) << setw(30) << cities[i].name << " |" << setw(13) << to_string(cities[i].population) << " |"
            << setw(15) << to_string(cities[i].area) << " |" << endl;
    }
    file << "--------------------------------------------------------------" << endl << endl;
    file.close();
}

void writeResultToFile(City cities[], string fileName, int res_size) {
    ofstream file;
    file.open(fileName, ios::app);
    file << setw(39) << "Rezultatai" << endl
        << "---------------------------------------------------------------------------------------" << endl
        << setw(5) << "Nr. |" << setw(30) << "Name |" << setw(15) << "Population |" << setw(17) << "Area |" << setw(19) << "Teksto rezultatas   |" << endl
        << "---------------------------------------------------------------------------------------" << endl;

    for (int i = 0; i < res_size; i++)
    {
        file << setw(5) << to_string(i+1) << setw(30) << cities[i].name << " |" << setw(13) << to_string(cities[i].population) << " |"
            << setw(15) << to_string(cities[i].area) << " |" << setw(17) << (cities[i].res) << " |" << endl;
    }
    file << "---------------------------------------------------------------------------------------" << endl;
    file.close();
}


__device__ void gpu_strcpy(char* dest, const char* src) {
	int i = 0;
	do {
		dest[i] = src[i];
	} while (src[i++] != 0);
}

__device__ void gpu_string(char* dest, const char* src) {
	dest[0] = src[0];
	dest[1] = src[1];
	dest[2] = '<';
	dest[3] = '1';
	dest[4] = '0';
	dest[5] = '0';
}

__global__ void gpu_func(City* device_cities, City* device_results, int* device_array_size, int* device_slice_size, int* device_result_count) {
    // Compute start index
    unsigned long start_index =* device_slice_size * threadIdx.x;
    unsigned long end_index;

    if (threadIdx.x == blockDim.x - 1) {
        end_index = *device_array_size;
    }
    else {
        end_index = *device_slice_size * (threadIdx.x + 1);
    }

    auto fp_sum = 0;

    for (int i = start_index; i < end_index; i++) {
        double population = device_cities[i].population;
        double area = device_cities[i].area;
        double density = population / area;
        if (density > 100) {
            City city;
            gpu_strcpy(city.name, device_cities[i].name);
            city.population = device_cities[i].population;
            city.area = device_cities[i].area;
            gpu_string(city.res, device_cities[i].name);

            // Inserting into results array
            int index = atomicAdd(device_result_count, 1);
            device_results[index] = city;
		}
	}
}

const int SIZE = 256;
const int BLOCKS = 2;
const int THREADS = 4;

int main()
{
	vector<City> data;
    readCitiesFile(&data);
    City* cities = &data[0];
    City results[SIZE];

    int slice_size = SIZE / THREADS;
    int result_count = 0;

    City* device_cities;
    City* device_results;
    int* device_array_size;
    int* device_slice_size;
    int* device_result_count;


    // Allocate memory on GPU
    cudaMalloc((void**)&device_cities, SIZE * sizeof(City));
    cudaMalloc((void**)&device_results, SIZE * sizeof(City));
    cudaMalloc((void**)&device_array_size, sizeof(int));
    cudaMalloc((void**)&device_slice_size, sizeof(int));
    cudaMalloc((void**)&device_result_count, sizeof(int));

    //Funkcijos, vykdomos GPU ir kviečiamos iš GPU
    //cudaMemcpyHostToHost iš CPU į CPU
    //cudaMemcpyHostToDevice iš CPU į GPU
    //cudaMemcpyDeviceToHost iš GPU į CPU
    //cudaMemcpyDeviceToDevice iš GPU į GPU

    // Copy data from CPU to GPU
    cudaMemcpy(device_cities, cities, SIZE * sizeof(City), cudaMemcpyHostToDevice);
    cudaMemcpy(device_array_size, &SIZE, sizeof(int), cudaMemcpyHostToDevice);
    cudaMemcpy(device_slice_size, &slice_size, sizeof(int), cudaMemcpyHostToDevice);
    cudaMemcpy(device_result_count, &result_count, sizeof(int), cudaMemcpyHostToDevice);

    gpu_func<<<BLOCKS, THREADS>>>(device_cities, device_results, device_array_size, device_slice_size, device_result_count);

    // Blocks CPU code until GPU code is done
    cudaDeviceSynchronize();

    cudaMemcpy(&results, device_results, SIZE * sizeof(City), cudaMemcpyDeviceToHost);
    int RES_SIZE = 0;
    cudaMemcpy(&RES_SIZE, device_result_count, sizeof(int), cudaMemcpyDeviceToHost);

    writeListToFile(data, resultFile);
    
    writeResultToFile(results, resultFile, RES_SIZE);

    // Free memory on 
    cudaFree(device_cities);
    cudaFree(device_results);
    cudaFree(device_array_size);
    cudaFree(device_slice_size);
    cudaFree(device_result_count);

    return 0;
}