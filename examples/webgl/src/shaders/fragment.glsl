precision mediump float;

varying float v_color;

void main() {
   gl_FragColor = vec4(0, v_color, v_color, 1);
}
