import numpy as np
import matplotlib.pyplot as plt
import time

# Sudarome tikslo funkciją
def tikslas(miesto_parduotuves, naujos_parduotuves):
    kainos = []
    for parduotuve in naujos_parduotuves:
        kaina = 0
        for miesto_parduotuve in miesto_parduotuves:
            atstumas = np.exp(-0.3 * ((parduotuve[0] - miesto_parduotuve[0]) ** 2 + (parduotuve[1] - miesto_parduotuve[1]) ** 2))
            kaina += atstumas
        vietos_kaina = ((parduotuve[0] ** 4 + parduotuve[1] ** 4) / 1000) + (np.sin(parduotuve[0]) + np.cos(parduotuve[1]) / 5) + 0.4
        kainos.append(kaina + vietos_kaina)
    return np.sum(kainos)

# Inicializuojame pradinius duomenis
np.random.seed(0)
n = 50  # Parduotuvių skaičius mieste
m = 50  # Kiek naujų parduotuvių pastatyti

# Sugeneruojame pradines parduotuves ir pridėtas mieste
miesto_parduotuves = np.random.rand(n, 2) * 20 - 10  # Koordinatės nuo -10 iki 10
pridetos_parduotuves = np.random.rand(m, 2) * 20 - 10

# Taikome gradientinį nusileidimą tikslo funkcijai
rate = 0.1
num_iterations = 100

start_time = time.time() 
tikslo_funkcijos_reiksmes = []
for i in range(num_iterations):
    grad = np.zeros((m, 2))
    for j in range(m):
        for k in range(2):
            perturbation = np.zeros((m, 2))
            perturbation[j][k] = 1e-10
            f_plus = tikslas(miesto_parduotuves, pridetos_parduotuves + perturbation)
            f_minus = tikslas(miesto_parduotuves, pridetos_parduotuves - perturbation)
            grad[j][k] = (f_plus - f_minus) / (2 * 1e-10)
    pridetos_parduotuves -= rate * grad
    tikslo_funkcijos_reiksmes.append(tikslas(miesto_parduotuves, pridetos_parduotuves))

end_time = time.time()  # End measuring time
runtime = end_time - start_time

# Pavaizduojame rezultatus
plt.figure(figsize=(8, 8))
plt.scatter(miesto_parduotuves[:, 0], miesto_parduotuves[:, 1], c='r', marker='o', label='Miesto parduotuvės')
plt.scatter(pridetos_parduotuves[:, 0], pridetos_parduotuves[:, 1], c='b', marker='x', label='Naujos parduotuvės')
plt.xlabel('X')
plt.ylabel('Y')
plt.legend(loc='best')
plt.title('Parduotuvių vietos')
plt.grid(True)
plt.show()


print("Iteracijų skaičius")
print(num_iterations)

# Tikslo funkcijos priklausomybės nuo iteracijų skaičiaus grafikas
plt.figure()
plt.plot(range(num_iterations), tikslo_funkcijos_reiksmes, marker='o')
plt.xlabel('Iteracijos')
plt.ylabel('Tikslo funkcijos reikšmė')
plt.title('Tikslo funkcijos priklausomybė nuo iteracijų skaičiaus')
plt.grid(True)
plt.show()

print(f"Tikslo funkcijos reiksme: {tikslo_funkcijos_reiksmes[-1]}")
print(f"Total runtime: {runtime} seconds")

# Ataskaitoje pateiksime pradinę ir gautą taškų konfigūraciją, tikslo funkcijos aprašymą, taikyto metodo pavadinimą ir parametrus,
# iteracijų skaičių, iteracijų pabaigos sąlygas ir tikslo funkcijos priklausomybės nuo iteracijų skaičiaus grafiką.