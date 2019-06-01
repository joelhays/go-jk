#version 410

in vec2 TexCoord;

uniform vec3 objectColor;

uniform sampler2D objectTexture;

out vec4 frag_color;

void main() {
    vec3 result = objectColor * vec3(texture(objectTexture, TexCoord));
    frag_color = vec4(result, 1.0f);
}