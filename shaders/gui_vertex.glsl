#version 410
layout (location = 0) in vec3 position;
layout (location = 2) in vec3 uv;

out vec2 TexCoord;

uniform mat4 model;

void main() {
    gl_Position = model * vec4(position, 1.0f);
    TexCoord = vec2(uv.x, 1.0 - uv.y);
}