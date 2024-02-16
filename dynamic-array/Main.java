// DynamicArrayExample.java
public class Main {
    public static void main(String[] args) {
        DynamicArray dynArray = new DynamicArray();

        // Add elements to the array
        dynArray.push(1);
        dynArray.push(2);
        dynArray.push(3);
        dynArray.print(); // Output: [1, 2, 3]

        // Remove element at index 1
        dynArray.remove(1);
        dynArray.print(); // Output: [1, 3]

        // Pop last element
        dynArray.pop();
        dynArray.print(); // Output: [1]

        // Check if the array is empty
        System.out.println("Is the array empty? " + dynArray.isEmpty()); // Output: false

        // Pop again
        dynArray.pop();
        dynArray.print(); // Output: []
        
        // Check if the array is empty again
        System.out.println("Is the array empty? " + dynArray.isEmpty()); // Output: true
    }
}
