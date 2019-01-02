import com.google.common.collect.ComparisonChain;
import lombok.AllArgsConstructor;
import lombok.EqualsAndHashCode;
import lombok.ToString;

import javax.annotation.Nullable;
import java.io.PrintStream;
import java.nio.file.Files;
import java.nio.file.Paths;
import java.util.ArrayList;
import java.util.Collections;
import java.util.List;

public class Main {
    private PrintStream out = System.out;
    private List<Cart> carts = new ArrayList<>();
    private char[][] map;
    private int height;
    private int width;

    public static void main(String[] args) throws Exception {
        List<String> input = Files.readAllLines(Paths.get("input.txt"));
        new Main().run(input);
    }

    private void run(List<String> input) {
        init(input);
        findCarts();
        fillCartPositions();
        print();
        for (int i = 0; i < 200000; i++) {
            move();
//            print();
        }
    }

    private void fillCartPositions() {
        carts.forEach(c -> {
            switch (map[c.x][c.y]) {
                case 'v':
                case '^':
                    map[c.x][c.y] = '|';
                    break;
                case '>':
                case '<':
                    map[c.x][c.y] = '-';
                    break;
                default:
                    throw new RuntimeException();
            }

        });
    }

    private void move() {
        Collections.sort(carts);
        carts.stream().filter(c -> !c.crashed).forEach(c -> {
            int newx = c.x + c.speed.x;
            int newy = c.y + c.speed.y;
            Cart crashed = getCart(newx, newy);

            if (crashed != null) {
                System.out.println("crash:" + newx + "/" + newy);
                c.crashed = true;
                crashed.crashed = true;
            }

            c.x = newx;
            c.y = newy;


            if ("/\\".contains("" + map[c.x][c.y])) {
                c.speed.rotate(map[c.x][c.y]);
            } else if (map[c.x][c.y] == '+') {
                c.turn();
            }
        });

        if (carts.stream().filter(c -> !c.crashed).count() == 1) {
            out.println(carts.stream().filter(c -> !c.crashed).findFirst().orElse(null));
            System.exit(0);
        }
    }

    private void findCarts() {
        carts = new ArrayList<>();
        for (int x = 0; x < width; x++) {
            for (int y = 0; y < height; y++) {
                if ("<>^v".contains("" + map[x][y])) {
                    carts.add(new Cart(x, y, map[x][y]));
                }
            }
        }
    }

    private void init(List<String> input) {
        height = input.size();
        width = input.get(0).length();
        map = new char[width][];
        for (var y = 0; y < height; y++) {
            String line = input.get(y);
            for (var x = 0; x < width; x++) {
                if (y == 0) {
                    map[x] = new char[height];
                }
                char c = line.charAt(x);
                map[x][y] = c;
            }
        }
    }

    private void print() {
        for (var y = 0; y < map[0].length; y++) {
            for (var x = 0; x < map.length; x++) {
                Cart c = getCart(x, y);
                if (c != null) {
                    out.print(c.speed);
                } else {
                    out.print(map[x][y]);
                }
            }
            out.println();
        }
        out.println();
    }

    @Nullable
    private Cart getCart(int x, int y) {
        return carts.stream()
                .filter(c -> !c.crashed)
                .filter(c -> c.x == x && c.y == y)
                .findFirst().orElse(null);
    }
}

@SuppressWarnings("SuspiciousNameCombination")
@AllArgsConstructor
class Vec {
    int x;
    int y;

    void rotate(char c) {
        if (c == '/') {
            int tmp = x;
            x = -y;
            y = -tmp;
        } else if (c == '\\') {
            int tmp = x;
            x = y;
            y = tmp;
        }
    }

    void rotateRight() {
        int tmp = x;
        x = -y;
        y = tmp;
    }

    void rotateLeft() {
        int tmp = x;
        x = y;
        y = -tmp;
    }

    @Override
    public String toString() {
        if (x == 1) {
            return ">";
        } else if (x == -1) {
            return "<";
        } else if (y == 1) {
            return "v";
        } else {
            return "^";
        }
    }
}

@SuppressWarnings({"SuspiciousNameCombination", "NullableProblems"})
@ToString
@EqualsAndHashCode
class Cart implements Comparable<Cart> {
    public int x;
    public int y;
    public Vec speed;
    public Turn nextTurn = Turn.LEFT;
    public boolean crashed = false;

    enum Turn {
        LEFT,
        NOPE,
        RIGHT
    }

    Cart(int x, int y, char c) {
        this.x = x;
        this.y = y;
        switch (c) {
            case '>':
                speed = new Vec(1, 0);
                break;
            case '<':
                speed = new Vec(-1, 0);
                break;
            case '^':
                speed = new Vec(0, -1);
                break;
            case 'v':
                speed = new Vec(0, 1);
                break;
        }
    }

    @Override
    public int compareTo(Cart o) {
        return ComparisonChain.start()
                .compare(y, o.y)
                .compare(x, o.x)
                .result();
    }

    void turn() {
        switch (nextTurn) {
            case LEFT:
                speed.rotateLeft();
                nextTurn = Turn.NOPE;
                break;
            case NOPE:
                nextTurn = Turn.RIGHT;
                break;
            case RIGHT:
                speed.rotateRight();
                nextTurn = Turn.LEFT;
                break;
        }
    }
}
