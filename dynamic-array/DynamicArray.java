public class DynamicArray {
    private int capacity; // Capacity of the dynamic array
    private int size;     // Number of elements currently in the array
    private int[] array;  // Array to hold elements

    // Constructor to initialize the dynamic array
    public DynamicArray() {
        capacity = 1;      // Initial capacity
        size = 0;          // No elements initially
        array = new int[capacity]; // Initialize the array with the initial capacity
    }

    // Method to add an element to the end of the array
    public void push(int element) {
        // Check if the array is full
        if (size == capacity) {
            // If full, double the capacity
            resize(2 * capacity);
        }
        // Add the element to the end of the array
        array[size] = element;
        // Increment the size
        size++;
    }

    // Method to remove and return the last element from the array
    public int pop() {
        // Check if the array is empty
        if (size == 0) {
            throw new IllegalStateException("Cannot pop from an empty array");
        }
        // Get the last element
        int lastElement = array[size - 1];
        // Decrement the size
        size--;
        // Return the last element
        return lastElement;
    }

    // Method to remove an element from a specific index
    public void remove(int index) {
        // Check if index is valid
        if (index < 0 || index >= size) {
            throw new IndexOutOfBoundsException("Index is out of bounds");
        }

        // Shift elements to the left to fill the gap
        for (int i = index; i < size - 1; i++) {
            array[i] = array[i + 1];
        }
        // Decrement the size
        size--;
    }

    // Method to resize the array to the new capacity
    private void resize(int newCapacity) {
        // Create a new array with the new capacity
        int[] newArray = new int[newCapacity];
        // Copy existing elements to the new array
        for (int i = 0; i < size; i++) {
            newArray[i] = array[i];
        }
        // Update the capacity and array reference
        capacity = newCapacity;
        array = newArray;
    }

    // Method to check if the array is empty
    public boolean isEmpty() {
        return size == 0;
    }

    // Method to print the elements of the array
    public void print() {
        System.out.print("[");
        for (int i = 0; i < size; i++) {
            System.out.print(array[i]);
            if (i < size - 1) {
                System.out.print(", ");
            }
        }
        System.out.println("]");
    }
}
