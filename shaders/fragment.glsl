#version 410

in vec3 FragPos;
in vec3 Normal;
in vec2 TexCoord;
in float LightIntensity;

uniform vec3 lightPos;
uniform vec3 viewPos;
uniform vec3 objectColor;
uniform vec3 lightColor;

uniform sampler2D objectTexture;

out vec4 frag_color;
void main() {
    // ambient
    float ambientStrength = 0.1f;
    vec3 ambient = ambientStrength * lightColor;

    // diffuse
    vec3 norm = normalize(Normal);
    vec3 lightDirection = normalize(lightPos - FragPos);
    float diff = max(dot(norm, lightDirection), 0.0);
    vec3 diffuse = diff * lightColor;

    // specular
    float specularStrength = 0.5f;
    vec3 viewDirection = normalize(viewPos - FragPos);
    vec3 reflectDirection = reflect(-lightDirection, norm);
    float spec = pow(max(dot(viewDirection, reflectDirection), 0.0), 32);
    vec3 specular = specularStrength * spec* lightColor;

    vec3 result = (ambient + diffuse + specular) * objectColor * vec3(texture(objectTexture, TexCoord));
    frag_color = vec4(result, 1.0f);

//    float strength = LightIntensity / 100.0f;
//    vec3 texColor = vec3(texture(objectTexture, TexCoord));
//    vec3 objColor = vec3(strength, strength, strength);
//    vec3 color = objColor + texColor;
//    color *= .35;
//    frag_color = vec4(color, 1.0f);

//    vec3 texColor = vec3(texture(objectTexture, TexCoord));
//    frag_color = vec4(texColor, 1.0f);
//    frag_color = vec4(1f, 1f, 1f, 1f);
}