import numpy as np
import matplotlib.pyplot as plt
import time
from multiprocessing import Pool

# Sudarome tikslo funkciją
def target(city_shops, new_shops):
    prices = []
    for new_shop in new_shops:
        price = 0
        for city_shop in city_shops:
            distance = np.exp(-0.3 * ((new_shop[0] - city_shop[0]) ** 2 + (new_shop[1] - city_shop[1]) ** 2))
            price += distance
        place_price = ((new_shop[0] ** 4 + new_shop[1] ** 4) / 1000) + (np.sin(new_shop[0]) + np.cos(new_shop[1]) / 5) + 0.4
        prices.append(price + place_price)
    return np.sum(prices)

def calculate_gradient(args):
    j, k, m, city_shops, addded_shops = args
    perturbation = np.zeros((m, 2))
    perturbation[j][k] = 1e-10
    f_plus = target(city_shops, addded_shops + perturbation)
    f_minus = target(city_shops, addded_shops - perturbation)
    return (f_plus - f_minus) / (2 * 1e-10)

def plot_shops(city_shops, added_shops):
    plt.figure(figsize=(8, 8))
    plt.scatter(city_shops[:, 0], city_shops[:, 1], c='r', marker='o', label='Miesto parduotuvės')
    plt.scatter(added_shops[:, 0], added_shops[:, 1], c='b', marker='x', label='Naujos parduotuvės')
    plt.xlabel('X')
    plt.ylabel('Y')
    plt.legend(loc='best')
    plt.title('Parduotuvių vietos')
    plt.grid(True)
    plt.show()

def plot_target_function(target_function, num_iterations):
    plt.figure()
    plt.plot(range(num_iterations), target_function, marker='o')
    plt.xlabel('Iteracijos')
    plt.ylabel('Tikslo funkcijos reikšmė')
    plt.title('Tikslo funkcijos priklausomybė nuo iteracijų skaičiaus')
    plt.grid(True)
    plt.show()

def plot_runtime(runtimes, dots, iterations):
    plt.figure()
    plt.plot(range(1, len(runtimes) + 1), runtimes, marker='o')
    plt.xlabel('Procesų skaičius')
    plt.ylabel('Laikas, s')
    plt.suptitle('Laiko priklausomybė nuo procesų skaičiaus')
    plt.title('Taškai: ' + str(dots) + ', Iteracijos: ' + str(iterations))
    plt.grid(True)
    plt.show()

# Inicializuojame pradinius duomenis
np.random.seed(0)
n = 10  # Parduotuvių skaičius mieste
m = 80  # Kiek naujų parduotuvių pastatyti
num_iterations = 128
max_processes = 8
runtime = []
rate = 0.1

# Sugeneruojame pradines parduotuves ir pridėtas mieste
city_shops = np.random.rand(n, 2) * 20 - 10  # Koordinatės nuo -10 iki 10
added_shops = np.random.rand(m, 2) * 20 - 10

for processes in range(1, max_processes + 1):
    if __name__ == '__main__':
        start_time = time.time() 
        target_function = []
        with Pool(processes) as p:
            for i in range(num_iterations):
                grad = np.array(p.map(calculate_gradient, [(j, k, m, city_shops, added_shops) for j in range(m) for k in range(2)]))
                grad = grad.reshape(m, 2)
                added_shops -= rate * grad
                target_function.append(target(city_shops, added_shops))
        end_time = time.time()  # End measuring time
        runtime.append(end_time - start_time)

        # Print results for each iteration within the loop
        print(f"Iteracijų skaičius: {num_iterations}")
        print(f"Tikslo funkcijos reiksme: {target_function[-1]}")
        print(f"Total runtime: {runtime[-1]} seconds\n")

if __name__ == '__main__':
    plot_runtime(runtime, m, num_iterations)