#version 330 core
// Input vertex data, different for all executions of this shader.
in vec4 pv;
void main(){

	gl_Position = pv;
}
