При реализации приложения использованы 4 структуры:
1. ```go
   type Message struct {
   Token  string
   FileID string
   Data   string
   }
   ```
2. ```go
   type Cache struct {
   mu          sync.Mutex           // используется при добавлении значений в map
   mapMessages map[string][]Message
   }
   ```
3. ```go
   type Worker struct {
   cache        *Cache              // Worker берет []Message для записи в файл
   interval     time.Duration       // устанавливает периодичность работы Worker-а
   wg           sync.WaitGroup      // необходимо дождаться завершения работы всех Worker-ов при graceful shutdown
   shutdownChan chan struct{}       // при получении сигнала из канала запускается очистка кэша
   }
    ```
4. ```go
   type WorkerFactory struct {
   mu              sync.Mutex      // используется при добавлении Worker-ов в массив и для итерации по массиву при вызове Stop()
   workers         []*Worker       // хранит массив Worker-ов для удобства управления количетсвом и отправки сигналов
   cache           *Cache          // при создании нового Worker-а передается ему в параметрах
   interval        time.Duration   // при создании нового Worker-а передается ему в параметрах
   wg              sync.WaitGroup  // используется для отслеживания созданных новых Worker-ов
   WorkersLoadChan chan struct{}   // передает сигнал на создание нового Worker-а, если в кэше еще есть данные
   shutdownChan    chan struct{}   // при получении сигнала запускается Stop(), при котором у каждого Worker-а вызывается этот же метод
   maxCountWorkers int             // показывает заданное максимально количество Worker-ов
   }
    ```



