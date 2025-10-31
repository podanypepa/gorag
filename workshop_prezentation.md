# Postav si vlastní AI: open source modely na tvém stroji

1. Úvod do workshopu
	•	Krátké představení lektora a účastníků, vysvětlení cíle workshopu.
	•	Cíl: pochopit, jak fungují open-source LLM a jak je rozjet lokálně.

2. Proč mít vlastní AI?
	•	Výhody lokálního provozu: nezávislost na cloudu, kontrola nad daty, žádné API limity.
	•	Reálné scénáře využití (AI asistent, analýza textů, výzkum, prototypování).

3. Typy modelů: uzavřené vs. otevřené
	•	Srovnání OpenAI, Anthropic, Google vs. open-source projekty (Meta, Mistral, Hugging Face).
	•	Výhody a nevýhody open-source přístupu.

4. Jak fungují LLM pod kapotou
	•	Stručný přehled: tokenizace, trénink, inferenční fáze.
	•	Představa o tom, co znamená 7B, 13B, 70B parametrů.
	•	Hardware požadavky (GPU, RAM, VRAM, CPU fallback).

5. Ekosystém open-source modelů
	•	Přehled platforem: Hugging Face, Ollama, LM Studio, vLLM, llama.cpp.
	•	Kde modely hledat, jak je hodnotit (MMLU, Benchmarks, community ratingy).

6. Výběr modelu pro osobní použití
	•	Kritéria: výkon vs. velikost, jazyková podpora, úkoly (chat, code, summarization…).
	•	Příklady: Mistral, Llama 3, Phi-3, Gemma, Starling, Zephyr.


💻 2. ČÁST – Praxe: Spusť si vlastní AI lokálně

1. Olama
    •	Ukázka, jak nainstalovat a spustit Ollamu (Windows/macOS/Linux).
    •	Vysvětlení, co Ollama dělá: lokální runtime pro LLM, správa modelů.

2. Vlastni RAG
    •	Ukázka, jak nastavit Retrieval-Augmented Generation (RAG) s lokálním modelem.

4. Vlastni model nebo vlastní data
    •	Jak si vytvořit vlastní model (modifikace parametrů, přidání system promptu).
    •	Zmínka o fine-tuning vs. LoRA vs. RAG.
    •	Ukázka jednoduchého custom modelu v Modelfile.

## 1. Úvod do workshopu

Velké jazykové modely (LLM = - Large Language Models) jsou hluboké neuronové sítě
(většinou architektury transformer) trénované na obrovských množstvích textových dat,
s cílem rozumět a generovat lidský text.  ￼
Modely jako GPT‑4, Claude nebo série Llama 3 jsou příklady, ale většina z nich není „otevřená“.

Klíčové pojmy
	•	Parametr – jeden váhový koeficient v modelu; model může mít miliardy (např. „13B“ znamená ~13 miliard parametrů).
	•	Transformér (Transformer) – architektura modelu, která umožňuje tzv. attention mechanismus.
	•	Tokenizace – převod vstupního textu na tokeny (menší jednotky, např. slova/posloupnosti znaků) pro model.
	•	Inference – fáze, kdy model „běží“ v reálném čase a generuje odpovědi; na rozdíl od tréninku.
	•	Fondový model (foundation model) – velký univerzální model, který lze dále upravovat (fine-tuning, LoRA, RAG).

## 2. Proč mít vlastní AI?

1. Soukromí a kontrola dat
    •	Žádná data neopouští tvůj počítač nebo lokální síť.
    •	Ideální pro citlivé projekty (interní dokumenty, kód, zákaznická data).
    •	Vyhneš se riziku úniku nebo zneužití informací přes cloudové služby.

2. Offline provoz
    •	Funguje i bez internetu.
    •	Užitečné pro vzdálené lokace, embedded zařízení nebo experimenty bez připojení.

3. Náklady a nezávislost
    •	Nemusíš platit za API volání (OpenAI, Anthropic, apod.).
    •	Nejsi závislý na změnách licencí, limitů nebo výpadcích poskytovatele.

4. Možnost ladění a rozšíření
    •	Můžeš upravovat model, přidávat vlastní data (fine-tuning, RAG).
    •	Experimenty s architekturami (Llama, Mistral, Gemma, Phi atd.).
    •	Plná kontrola nad prompt engineeringem a pipeline.

## 3. Typy modelů: uzavřené vs. otevřené

### 🧠 Co je to open-source LLM

LLM (Large Language Model) je umělá inteligence natrénovaná na obrovském množství textů,
která dokáže rozumět přirozenému jazyku, odpovídat, generovat text, překládat, psát kód a mnoho dalšího.

🪟 Open-source LLM = otevřený jazykový model

🗝️ „Open-source“ znamená, že model je veřejně dostupný – můžeš si stáhnout jeho váhy, kód i architekturu
a používat ho bez nutnosti připojení ke cloudu.

📦 Co přesně bývá otevřené
    •	Modelové váhy (weights) – samotné naučené parametry modelu.
    •	Architektura – popis, jak model funguje (počet vrstev, attention, embeddingy…).
    •	Kód pro inference / trénování – skripty, které umožňují model spustit nebo učit.
    •	Licence – právní rámec určující, co s modelem smíš dělat (např. MIT, Apache-2.0, Meta License).

### 🔒 Uzavřené (proprietární) modely

Příklady: GPT-4 (OpenAI), Claude 3 (Anthropic), Gemini (Google), Copilot (Microsoft)

✅ Výhody:
    •	Nejvyšší kvalita výstupů (velké množství tréninkových dat, optimalizace).
    •	Spolehlivost, bezpečnost, stabilita.
    •	Připravené integrace (API, pluginy, Copiloty).
    •	Žádné starosti s hardwarem — běží v cloudu.

⚠️ Nevýhody:
    •	Uzavřený kód a model — nevíš, jak byl natrénován.
    •	Omezení licencí a podmínek (např. zákaz určitého použití).
    •	Data putují do cloudu (→ otázky soukromí).
    •	Náklady podle počtu tokenů / API volání.
    •	Nelze model přetrénovat nebo přizpůsobit.

### 🧠 Otevřené (open-source) modely

Příklady: Llama 3 (Meta), Mistral, Gemma (Google), Phi-3 (Microsoft), Falcon, DeepSeek

✅ Výhody:
    •	Plná kontrola — můžeš model spustit, upravit, kvantizovat, nebo trénovat dál.
    •	Běží lokálně nebo v privátním cloudu → žádný únik dat.
    •	Můžeš integrovat do vlastních aplikací bez omezení.
    •	Komunita vytváří vylepšené varianty (instruct, RAG, agenti).
    •	Většina licencí je velmi otevřená (Apache-2.0, MIT, Meta License).

⚠️ Nevýhody:
    •	Výkon závisí na tvém hardwaru (RAM, GPU).
    •	Kvalita může být nižší než u top uzavřených modelů.
    •	Vyžaduje znalosti ohledně nasazení a optimalizace.
    •	Menší modely = menší „chytrost“, větší = náročné na zdroje.

##  4. Jak fungují LLM pod kapotou

### 🧩 1️⃣ Co LLM vlastně dělá

LLM (Large Language Model) je neuronová síť, která se naučila předpovídat další slovo (token) na základě předchozích.

Příklad:
„Pes běžel přes …“ → model se snaží odhadnout, že další token bude „louku“.

🔡 2️⃣ Tokenizace – jak se text mění na čísla

Text	        Tokeny	            Počet tokenů
„Ahoj světe!“	[Ahoj, svě, te, !]	4
„Hello world!“	[Hello, world, !]	3

Potom se tokeny převedou na číselné vektory (embeddingy), se kterými model počítá.

🧩 Každý token je jako souřadnice významu ve vícerozměrném prostoru.

### 🧱 3️⃣ Trénink – jak se model učí jazyk

Model se učí předpovídat další token podle předchozích.
Trénuje se na obrovských textech (např. internet, knihy, kód, Wikipedie).
Optimalizuje miliony/miliardy parametrů (čísel, která určují, „co si pamatuje“).
Používá GPU farmy – trénink trvá týdny nebo měsíce.

💬 Výsledek: model si „osvojí“ statistické vzorce jazyka – syntax, logiku, fakta i styl.

### ⚙️ 4️⃣ Inferenční fáze – jak model odpovídá

To je to, co se děje při každém chatu:
•	Zadáš prompt → text se tokenizuje.
•	Model spočítá pravděpodobnosti pro každý možný další token.
•	Vybere nejvhodnější (nebo nejpravděpodobnější).
•	Tento token přidá do vstupu a proces opakuje.

➡️ Takto vzniká odpověď token po tokenu – např. slovo po slově.

### 🧮 5️⃣ Co znamená 7B, 13B, 70B parametrů

„B“ = miliardy parametrů (weights, tj. čísla uložená v modelu).
Čím víc parametrů → tím větší kapacita modelu učit se vzorce.

Model	Parametrů	    Velikost        (FP16)	Typické využití
7B	    7 miliard	    13 GB	        Lokální běh, rychlé odpovědi
13B	    13 miliard	    26 GB	        Vyvážený výkon
70B	    70 miliard	    140 GB	        Výzkum, servery, top kvalita

💡 Menší modely jsou rychlejší, ale méně přesné; větší jsou chytřejší, ale náročnější.

### 🧰 6️⃣ Hardware požadavky

💻 GPU (doporučeno)
Nejlepší pro inference (paralelní výpočty).
Např. NVIDIA RTX 4070/4090 zvládne 7B–13B modely.
U velkých modelů (70B+) se používají více-GPU servery.

🧠 RAM
Kvantizované modely (např. GGUF Q4) se vejdou do 8–16 GB RAM.
Nekvantizované modely vyžadují desítky GB RAM.

🎮 VRAM
Každý GB VRAM ≈ ~1 miliarda parametrů (zjednodušeně).
7B → 8–12 GB VRAM
13B → 16–24 GB VRAM
70B → 80–140 GB VRAM (server only)

⚙️ CPU fallback
Lze spustit i bez GPU (Ollama, llama.cpp, LM Studio).
Výrazně pomalejší, ale pro menší modely použitelné.
Užitečné pro testování nebo RAG úlohy bez interaktivního chatu.

## 5. Ekosystém open-source modelů

🌍 1️⃣ Co znamená „ekosystém“ open-source modelů

Open-source AI komunita roste neuvěřitelně rychle – stovky týmů po celém světě vyvíjejí otevřené LLM.
To znamená:
    •	Každý si může model stáhnout, spustit, upravit nebo natrénovat.
    •	Existují stovky variant (7B, 13B, 70B, kvantizace, fine-tuny).
    •	Vznikl celý ekosystém platforem a nástrojů, které to zpřehledňují.

### 🧰 2️⃣ Přehled hlavních platforem a nástrojů

🧡 Hugging Face
    •	Největší katalog AI modelů a datasetů.
    •	Hostuje tisíce LLM, nabízí testování v prohlížeči, API, dokumentaci.
    •	Vyhledávání, porovnávání, sdílení modelů.

🦙 Ollama
    •	Nástroj pro snadné spuštění modelů lokálně (ollama run mistral).
    •	Podporuje Llama, Mistral, Phi, Gemma a další.
    •	Lokální inference a integrace do aplikací.

💻 LM Studio
    •	GUI aplikace pro Windows/macOS/Linux.
    •	Umožňuje stahovat a spouštět modely vizuálně, podobně jako ChatGPT.
    •	Demo, výuka, rychlé testy modelů.

⚡ vLLM
    •	Výkonný inference engine vyvinutý výzkumníky z UC Berkeley.
    •	Optimalizovaný pro rychlost a škálování (API server, GPU cluster).
    •	Produkční nasazení, cloud nebo server.

🧮 llama.cpp
    •	C++ implementace pro běh modelů v kvantizované podobě (GGUF).
    •	Velmi efektivní – běží i na CPU.
    •	Lokální testování, embedded, offline použití.

### 📊 4️⃣ Jak modely hodnotit (benchmarks a metriky)

LLM se liší kvalitou, rychlostí i velikostí.
Existují standardizované benchmarks, které měří různé aspekty schopností:

#### Metrika, Co měří, Typ testů

MMLU (Massive Multitask Language Understanding)
•	Obecná znalost napříč 57 obory (věda, historie, právo, IT).
•	Výběr odpovědi z možností.

ARC, HellaSwag, TruthfulQA
•	Logické a faktické uvažování.
•	Testy z reálných scénářů.

GSM8K
•	Matematické úlohy pro 8. třídu.
•	Výpočty a krokové uvažování.

HumanEval
•	Programování – dokáže napsat funkci podle zadání?
•	Generování kódu.

## 6. Výběr modelu pro osobní použití

### 🔍 Kritéria pro výběr modelu

1. Výkon vs. velikost

Velikost modelu (např. 7B, 13B, 70B parametrů) silně ovlivňuje výpočetní nároky i kvalitu výstupu.
Menší modely → méně parametrů → rychlejší běh, menší hardwarové požadavky, ale mohou mít horší výkon ve složitých úlohách.
Větší modely → vyšší „kapacita“ pro učení jazykových vzorců, kontextu, kódu atp., ale potřebují silnější HW (GPU, více VRAM/RAM).

2. Jazyková podpora a doména úkolů

Zvaž, zda model podporuje jazyk, ve kterém budeš pracovat (např. čeština, slovenština, angličtina)
    → menší modely nebo starší verze mohou mít horší podporu menších jazyků.
Úkoly: odpovídání na dotazy („chat“), generování kódu, sumarizace textu, překlad, RAG nad vlastním datasetem.
    → Každý model má ve svém profilu, co zvládá dobře.

3. Hardwarové a provozní omezení

Jaký máš hardware: GPU výkonné, nebo jen CPU? Kolik VRAM a RAM máš? Lokální inference nebo cloud?
Lokální použití: menší modely (např. ~7B) jsou mnohem praktičtější pro běh na běžném PC/GPU či dokonce na CPU s kvantizací.
Pokud nasazuješ v produkci nebo máš server, můžeš zvolit větší modely.
Ujisti se, že model má otevřenou licenci (Apache-2.0, MIT, podobně) a že ti dovoluje ho používat způsobem, jak potřebuješ.
